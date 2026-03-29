package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	id            int
	nick          string
	password_hash string
}

type Credentials struct {
	Nick     string `json:"nick"`
	Password string `json:"password"`
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type Product struct {
	Id          int     `json:"id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	CategoryId  int     `json:"category_id"`
}

type Categorie struct {
	Id           int    `json:"id"`
	CategoryName string `json:"category_name"`
}
