package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware
	r := gin.Default()

	// Define a simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Define a simple POST endpoint
	r.POST("/echo", func(c *gin.Context) {
		var json struct {
			Message string `json:"message"`
		}
		
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"echo": json.Message,
		})
	})

	// Start server on port 8080
	r.Run(":8080")
}
