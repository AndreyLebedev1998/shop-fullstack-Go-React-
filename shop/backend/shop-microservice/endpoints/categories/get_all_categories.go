package categories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"shop/models"
)

func getAllCategories(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var categories []models.Categorie
	var ctx context.Context = context.Background()
	var query string = "SELECT * FROM categories"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		http.Error(w, "Error while querying the database", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var category models.Categorie
		err := rows.Scan(&category.Id, &category.CategoryName)
		if err != nil {
			fmt.Println("Error reading line")
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		categories = append(categories, category)
	}

	json.NewEncoder(w).Encode(categories)
}
