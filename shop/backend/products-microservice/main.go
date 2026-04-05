package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"products-microservice/cors"
	"products-microservice/endpoints/categories"
	"products-microservice/endpoints/products"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/redis/go-redis/v9"

	_ "products-microservice/docs"

	_ "github.com/lib/pq"
)

// @title Vanilla Go API
// @version 1.0
// @description prodcuts-microservice
// @host localhost:8090
// @BasePath /

func main() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	var db *sql.DB

	psqlInfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("✅ Подключено к PostgreSQL products-microservice")

	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:        os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		DB:          0,
		MaxRetries:  3,
		DialTimeout: 3 * time.Second,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	fmt.Println("✅ Подключено к Redis products-microservice")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🔥 Работает! Сервер перезапустился!"))
	})

	mux := http.NewServeMux()

	mux.Handle("/products", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.GetAllProductsByCategoryId(w, r, db, rdb)
	})))

	mux.Handle("/categories", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		categories.GetAllCategories(w, r, db, rdb)
	})))

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	http.ListenAndServe(":8080", mux)
}
