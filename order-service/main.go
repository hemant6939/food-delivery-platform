package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Order struct {
	ID             int       `json:"id"`
	CustomerName   string    `json:"customer_name"`
	RestaurantName string    `json:"restaurant_name"`
	Items          string    `json:"items"`
	TotalAmount    float64   `json:"total_amount"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

var db *sql.DB

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPass, dbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to PostgreSQL")

	r := gin.Default()
	r.Use(corsMiddleware())
	r.POST("/api/orders", createOrder)
	r.GET("/api/orders", getOrders)
	r.GET("/health", health)
	r.Run(":8080")
}

func createOrder(c *gin.Context) {
	var order Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(
		"INSERT INTO orders (customer_name, restaurant_name, items, total_amount, status, created_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		order.CustomerName, order.RestaurantName, order.Items, order.TotalAmount, "PLACED", time.Now(),
	).Scan(&order.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order.Status = "PLACED"
	c.JSON(http.StatusCreated, order)
}

func getOrders(c *gin.Context) {
	rows, err := db.Query("SELECT id, customer_name, restaurant_name, items, total_amount, status, created_at FROM orders")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		rows.Scan(&o.ID, &o.CustomerName, &o.RestaurantName, &o.Items, &o.TotalAmount, &o.Status, &o.CreatedAt)
		orders = append(orders, o)
	}
	c.JSON(http.StatusOK, orders)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}