package orders

import (
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

	if len(order.OrderItems) == 0 {
		http.Error(w, "Order must contain items", http.StatusBadRequest)
		return
	}

	if order.Status == nil || *order.Status != "pending" {
		val := "pending"
		order.Status = &val
	}

	var ctx = r.Context()

	var queryOrder string = `INSERT INTO orders (user_id, email, phone, status, total_price)
						VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var queryOrderItems string = `INSERT INTO order_items (order_id, product_id, quantity, price)
									VALUES ($1, $2, $3, $4)`

	var newOrderId models.NewOrderId

	// Начало транзакции
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, queryOrder, order.UserId, order.Email, order.Phone, order.Status, order.TotalPrice).Scan(&newOrderId.Id)
	if err != nil {
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

	for _, item := range order.OrderItems {
		_, err := tx.ExecContext(ctx, queryOrderItems, newOrderId.Id, item.ProductId, item.Quantity, item.Price)
		if err != nil {
			http.Error(w, "Error creating order_items", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}
	// Конец транзакции

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newOrderId)
}
