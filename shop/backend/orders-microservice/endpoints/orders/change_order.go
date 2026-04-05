package orders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"orders-microservice/models"
)

func ChangeOrder(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order models.Order
	var newOrder models.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if len(order.OrderItems) == 0 {
		http.Error(w, "Order must contain items", http.StatusBadRequest)
		return
	}

	var orderItems []models.OrderItem = order.OrderItems
	var ctx = r.Context()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback()

	var deleteQuery string = `DELETE FROM order_items WHERE order_id = $1`
	var insertQuery string = `INSERT INTO order_items (order_id, product_id, quantity, price)
									VALUES ($1, $2, $3, $4)`

	res, err := tx.ExecContext(ctx, deleteQuery, order.Id)

	if err != nil {
		fmt.Println("DB update error:", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	deleteRows, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if deleteRows == 0 {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	for _, item := range orderItems {
		_, err := tx.ExecContext(ctx, insertQuery, order.Id, item.ProductId, item.Quantity, item.Price)
		if err != nil {
			http.Error(w, "Error creating order_items", http.StatusInternalServerError)
			return
		}
	}

	itemsQuery := `SELECT product_id, quantity, price FROM order_items WHERE order_id = $1`
	rows, err := tx.QueryContext(ctx, itemsQuery, order.Id)

	if err != nil {
		http.Error(w, "Error fetching order items", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ProductId, &item.Quantity, &item.Price); err != nil {
			http.Error(w, "Error scanning order items", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error reading order items", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}
	newOrder = order
	newOrder.OrderItems = items
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newOrder)
}
