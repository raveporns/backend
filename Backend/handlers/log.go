package handlers

import (
	"net/http"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

func GetLog(c *gin.Context) {
    rows, err := database.DB.Query("SELECT id,client_id,endpoint,method,timestamp FROM logs")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch log"})
        return
    }
    defer rows.Close()

    // Initialize เป็น slice ว่างแทนที่จะเป็น nil
    results := make([]models.Log, 0)
    for rows.Next() {
        var log models.Log
        if err := rows.Scan(&log.ID, &log.ClientID,&log.Endpoint,&log.Method,&log.Time); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Log scan error"})
            return
        }
        results = append(results, log)
    }

    c.JSON(http.StatusOK, results)
}
func GetLogByID(c *gin.Context) {
    client_id := c.Param("id")

    // ดึง log เพียงแค่รายการเดียวตาม id
    row := database.DB.QueryRow("SELECT id, client_id, endpoint, method, timestamp FROM logs WHERE client_id = $1", client_id)

    var log models.Log
    if err := row.Scan(&log.ID, &log.ClientID, &log.Endpoint, &log.Method, &log.Time); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
        return
    }

    c.JSON(http.StatusOK, log)
}

func PostLog(c *gin.Context) {
	var log models.Log

	// Bind JSON ที่ส่งมาจาก client
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// บันทึก log ลงฐานข้อมูล
	query := `INSERT INTO log (client_id, endpoint, method, time) VALUES (?, ?, ?, ?)`
	_, err := database.DB.Exec(query, log.ClientID, log.Endpoint, log.Method, log.Time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not insert log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Log saved"})
}
