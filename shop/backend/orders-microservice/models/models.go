package models

type Order struct {
	Id         int     `json:"id"`
	UserId     *int    `json:"user_id"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Status     *string `json:"status"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type NewOrder struct {
	UserId     *int    `json:"user_id"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Status     *string `json:"status"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type NewOrderId struct {
	Id int `json:"id"`
}
