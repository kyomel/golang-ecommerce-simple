package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"onlinetoko/handler"
	"onlinetoko/middleware"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v\n", err)
	}

	db, err := sql.Open("pgx", os.Getenv("DB_URI"))
	if err != nil {
		log.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	if _, err = migrate(db); err != nil {
		log.Printf("Failed to migrate database: %v\n", err)
		os.Exit(1)
	}

	r := gin.Default()

	r.GET("/api/v1/products", handler.ListProducts(db))
	r.GET("/api/v1/products/:id", handler.GetProducts(db))
	r.POST("/api/v1/checkout", handler.CheckoutOrder(db))

	r.POST("/api/v1/orders/:id/confirm", handler.ConfirmOrder(db))
	r.GET("/api/v1/orders/:id", handler.GetOrder(db))

	r.POST("/admin/products", middleware.AdminOnly(), handler.CreateProduct(db))
	r.PUT("/admin/products/:id", middleware.AdminOnly(), handler.UpdateProduct(db))
	r.DELETE("/admin/products/:id", middleware.AdminOnly(), handler.DeleteProduct(db))

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err = server.ListenAndServe(); err != nil {
		fmt.Printf("Failed to run server %v\n", err)
		os.Exit(1)
	}
}

func migrate(db *sql.DB) (sql.Result, error) {
	if db == nil {
		return nil, errors.New("connection is not available")
	}

	return db.Exec(`
	CREATE TABLE IF NOT EXISTS products (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price BIGINT NOT NULL,
		is_deleted BOOLEAN NOT NULL DEFAULT FALSE
	);

	CREATE TABLE IF NOT EXISTS orders (
		id VARCHAR(36) PRIMARY KEY,
		email VARCHAR(255) NOT NULL,
		address VARCHAR NOT NULL,
		passcode VARCHAR,
		paid_at TIMESTAMP,
		paid_bank VARCHAR(255),
		paid_account VARCHAR(255),
		grand_total BIGINT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS order_details (
		id VARCHAR(36) PRIMARY KEY,
		order_id VARCHAR(36) NOT NULL,
		product_id VARCHAR(36) NOT NULL,
		quantity INT NOT NULL,
		price BIGINT NOT NULL,
		total BIGINT NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id) ON UPDATE CASCADE ON DELETE RESTRICT,
		FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT
	);
	`)
}
