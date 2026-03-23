package products

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"shop/models"
	"strconv"
)

func GetAllProductsByCategoryId(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var category_id string
	var idStr string = r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		http.Error(w, "Invalid category_id", http.StatusBadRequest)
		return
	}

	if idStr == "" {
		http.Error(w, "category_id cannot be empty", http.StatusBadRequest)
		return
	}

	category_id = idStr

	var products []models.Product
	var query string = "SELECT * FROM products WHERE category_id = $1"
	var ctx context.Context = context.Background()

	rows, err := db.QueryContext(ctx, query, category_id)

	if err != nil {
		http.Error(w, "Error while querying the database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var product models.Product

		err := rows.Scan(&product.Id, &product.ProductName, &product.CategoryId, &product.Price)

		if err != nil {
			fmt.Println("Error reading line")
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		products = append(products, product)
	}

	json.NewEncoder(w).Encode(products)
}
