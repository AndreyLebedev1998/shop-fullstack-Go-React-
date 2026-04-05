package orders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"orders-microservice/models"
)

func GetOrdersByParametr(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var emailParametr string = r.URL.Query().Get("email")
	var phoneParamert string = r.URL.Query().Get("phone")
	var userIdParametr string = r.URL.Query().Get("user_id")
	var ctx = r.Context()
	var fullOrders []models.FullOrder
	ordersMap := make(map[int]*models.FullOrder)
	var query string = `SELECT orders.id as order_id, user_id, email, phone, status, total_price, created_at, order_items.id AS order_item_id, 
							product_id, quantity, order_items.price, product_name, category_id, category_name, image_url
							FROM orders 
							JOIN order_items ON orders.id = order_items.order_id
							JOIN products ON order_items.product_id = products.id
							JOIN categories ON products.category_id = categories.id`
	var sqlParamEmail = "WHERE email = $1"
	var sqlParamPhone = "WHERE phone = $1"
	var sqlParamUserId = "WHERE user_id = $1"

	if emailParametr != "" {

		rows, err := db.QueryContext(ctx, query+" "+sqlParamEmail, emailParametr)

		if err != nil {
			http.Error(w, "Error while querying the database", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var fullOrder models.FullOrder
			var product models.Products
			err := rows.Scan(&fullOrder.OrderId, &fullOrder.UserId, &fullOrder.Email, &fullOrder.Phone, &fullOrder.Status,
				&fullOrder.TotalPrice, &fullOrder.CreatedAt, &product.OrderItemId, &product.ProductId,
				&product.Quantity, &product.Price, &product.ProductName, &product.CategoryId, &product.CategoryName, &product.ImageUrl)
			if err != nil {
				fmt.Println("Error reading line")
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			if existingOrder, ok := ordersMap[fullOrder.OrderId]; ok {
				existingOrder.Products = append(existingOrder.Products, product)
			} else {
				fullOrder.Products = []models.Products{product}
				ordersMap[fullOrder.OrderId] = &fullOrder
			}
		}

		for _, order := range ordersMap {
			fullOrders = append(fullOrders, *order)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fullOrders)
	}

	if phoneParamert != "" {
		rows, err := db.QueryContext(ctx, query+" "+sqlParamPhone, emailParametr)

		if err != nil {
			http.Error(w, "Error while querying the database", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var fullOrder models.FullOrder
			var product models.Products
			err := rows.Scan(&fullOrder.OrderId, &fullOrder.UserId, &fullOrder.Email, &fullOrder.Phone, &fullOrder.Status,
				&fullOrder.TotalPrice, &fullOrder.CreatedAt, &product.OrderItemId, &product.ProductId,
				&product.Quantity, &product.Price, &product.ProductName, &product.CategoryId, &product.CategoryName)
			if err != nil {
				fmt.Println("Error reading line")
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			if existingOrder, ok := ordersMap[fullOrder.OrderId]; ok {
				existingOrder.Products = append(existingOrder.Products, product)
			} else {
				fullOrder.Products = []models.Products{product}
				ordersMap[fullOrder.OrderId] = &fullOrder
			}
		}

		for _, order := range ordersMap {
			fullOrders = append(fullOrders, *order)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fullOrders)
	}

	if userIdParametr != "" {
		rows, err := db.QueryContext(ctx, query+" "+sqlParamUserId, emailParametr)

		if err != nil {
			http.Error(w, "Error while querying the database", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var fullOrder models.FullOrder
			var product models.Products
			err := rows.Scan(&fullOrder.OrderId, &fullOrder.UserId, &fullOrder.Email, &fullOrder.Phone, &fullOrder.Status,
				&fullOrder.TotalPrice, &fullOrder.CreatedAt, &product.OrderItemId, &product.ProductId,
				&product.Quantity, &product.Price, &product.ProductName, &product.CategoryId, &product.CategoryName)
			if err != nil {
				fmt.Println("Error reading line")
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			if existingOrder, ok := ordersMap[fullOrder.OrderId]; ok {
				existingOrder.Products = append(existingOrder.Products, product)
			} else {
				fullOrder.Products = []models.Products{product}
				ordersMap[fullOrder.OrderId] = &fullOrder
			}
		}

		for _, order := range ordersMap {
			fullOrders = append(fullOrders, *order)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fullOrders)
	}
}
