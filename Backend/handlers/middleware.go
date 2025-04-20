package handlers

import (
	"project/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key required"})
			c.Abort()
			return
		}

		var clientID int
		err := database.DB.QueryRow("SELECT id FROM clients WHERE api_key = $1", apiKey).Scan(&clientID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		c.Set("client_id", clientID)
		c.Next()
	}
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, exists := c.Get("client_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		_, err := database.DB.Exec("INSERT INTO logs (client_id, endpoint, method) VALUES ($1, $2, $3)",
			clientID, c.Request.URL.Path, c.Request.Method)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not log request"})
			c.Abort()
			return
		}

		c.Next()
	}
}
