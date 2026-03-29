package orders

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"orders-microservice/models"
)

func CreateOrder(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order models.NewOrder

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var ctx = context.Background()
	var query string = `INSERT INTO orders (user_id, email, phone, status, total_price, created_at, updated_at)
						VALUES ($1, $2, $3, $4, $5, $6, &7) RETURNING id`

	if order.Status == nil || *order.Status != "pending" {
		val := "pending"
		order.Status = &val
	}

	var newOrderId models.NewOrderId

	err := db.QueryRowContext(ctx, query, order.UserId, order.Email, order.Phone, order.Status, order.TotalPrice, order.CreatedAt, order.UpdatedAt).Scan(&newOrderId)
	if err != nil {
		http.Error(w, "Error inserting into the database", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newOrderId)
}
