package create

import (
	"admin-microservice/models"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func CreateProduct(w http.ResponseWriter, r *http.Request, sql *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var ctx = r.Context()

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	fmt.Println("FORM:", r.Form)

	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	categoryIdStr := r.FormValue("category_id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		http.Error(w, "Invalid category_id", http.StatusBadRequest)
		return
	}

	nameStr := r.FormValue("name")
	if nameStr == "" {
		http.Error(w, "name is not defined", http.StatusBadRequest)
		return
	}

	var newProduct models.NewProduct = models.NewProduct{
		ProductName: nameStr,
		Price:       price,
		CategoryId:  categoryId,
		ImageUrl:    nil,
	}

	file, header, err := r.FormFile("image_url")
	if err != nil {
		if err != http.ErrMissingFile {
			http.Error(w, "Invalid file upload", http.StatusBadRequest)
			return
		}
	} else {
		defer file.Close()

		uploadDir := "./uploads-products-images"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		ext := filepath.Ext(header.Filename)
		safeName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), rand.Text())
		destPath := filepath.Join(uploadDir, safeName+ext)

		dst, err := os.Create(destPath)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		imageURL := fmt.Sprintf("/uploads-products-images/%s", filepath.Base(destPath))
		newProduct.ImageUrl = &imageURL
	}

	query := `INSERT INTO products (product_name, category_id, price, image_url)
			  VALUES ($1, $2, $3, $4) RETURNING id`

	var productId int
	err = sql.QueryRowContext(ctx, query, newProduct.ProductName, newProduct.CategoryId, newProduct.Price, newProduct.ImageUrl).Scan(&productId)

	if err != nil {
		http.Error(w, "Error saving product", http.StatusInternalServerError)
		return
	}

	var product = models.Product{
		Id:          productId,
		ProductName: newProduct.ProductName,
		CategoryId:  newProduct.CategoryId,
		Price:       newProduct.Price,
		ImageUrl:    newProduct.ImageUrl,
	}

	json.NewEncoder(w).Encode(product)
}
