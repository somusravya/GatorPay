package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// Get all transactions
	router.GET("/api/transactions", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"transactions": []gin.H{
				{
					"id":     1,
					"payer":  "Alice",
					"amount": 50.00,
					"date":   "2026-02-12",
				},
				{
					"id":     2,
					"payer":  "Bob",
					"amount": 30.00,
					"date":   "2026-02-11",
				},
			},
		})
	})

	// Create a new transaction
	router.POST("/api/transactions", func(c *gin.Context) {
		var transaction gin.H
		if err := c.BindJSON(&transaction); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{
			"id":        3,
			"payer":     transaction["payer"],
			"amount":    transaction["amount"],
			"date":      transaction["date"],
			"message":   "Transaction created successfully",
		})
	})

	// Get users
	router.GET("/api/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"users": []gin.H{
				{
					"id":    1,
					"name":  "Alice",
					"email": "alice@example.com",
				},
				{
					"id":    2,
					"name":  "Bob",
					"email": "bob@example.com",
				},
			},
		})
	})

	router.Run(":8080")
}
