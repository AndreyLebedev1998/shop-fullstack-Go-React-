package create

import (
	"admin-microservice/models"
	"database/sql"
	"encoding/json"
	"net/http"
)

func CreateCategory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var categoryName models.NewCategory

	if err := json.NewDecoder(r.Body).Decode(&categoryName); err != nil {
		http.Error(w, "category_name is not defined", http.StatusBadRequest)
		return
	}
}
