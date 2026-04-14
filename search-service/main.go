package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	esURL := os.Getenv("ES_URL")

	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/api/search", func(c *gin.Context) {
		query := c.Query("q")

		esQuery := fmt.Sprintf(`{
			"query": {
				"multi_match": {
					"query": "%s",
					"fields": ["customer_name", "restaurant_name", "items"]
				}
			}
		}`, query)

		req, _ := http.NewRequest("GET", esURL+"/orders/_search", strings.NewReader(esQuery))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.Data(200, "application/json", body)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Run(":8081")
}