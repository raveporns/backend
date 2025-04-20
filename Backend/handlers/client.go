package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

func generateAPIKey() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func RegisterClient(c *gin.Context) {
	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	client.APIKey = generateAPIKey()

	_, err := database.DB.Exec("INSERT INTO clients (name, email, api_key) VALUES ($1, $2, $3)", client.Name, client.Email, client.APIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Client registered", "api_key": client.APIKey})
}


func GetClient(c *gin.Context) {
	email := c.GetHeader("X-Email")

	var client models.Client
	row := database.DB.QueryRow("SELECT id, name, email, api_key FROM clients WHERE email = $1", email)
	err := row.Scan(&client.ID, &client.Name, &client.Email, &client.APIKey)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, client)
}

