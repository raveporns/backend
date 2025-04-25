package handlers

import (
	"database/sql"
	"net/http"
	"project/database"
	"project/models"
	"time"

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

func GetClicksSearch(c *gin.Context) {
	// นับตั้งแต่เที่ยงคืนวันนี้
	startOfDay := time.Now().Truncate(24 * time.Hour)

	var count int
	err := database.DB.QueryRow(`
	  SELECT COUNT(*) 
	  FROM logs 
	  WHERE timestamp >= $1
	`, startOfDay).Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			count = 0
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"clicksSearch": count})
}

func PostClickLog(c *gin.Context) {
    var log models.Log
    if err := c.ShouldBindJSON(&log); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // กรณีที่เวลาไม่ได้ถูกส่งมา, จะใช้เวลาปัจจุบัน
    if log.Time.IsZero() {
        log.Time = time.Now().UTC()
    }

    // บันทึก log ลงในฐานข้อมูล
    _, err := database.DB.Exec(
        `INSERT INTO logs (client_id, endpoint, method, timestamp)
         VALUES ($1, $2, $3, $4)`,
        log.ClientID, log.Endpoint, log.Method, log.Time,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not insert click log"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Click log saved"})
}

// GET /stats/link-clicks
func GetLinkClicks(c *gin.Context) {
    startOfDay := time.Now().Truncate(24 * time.Hour)

    var count int
    // นับเฉพาะ log ที่ endpoint เริ่มด้วย "/detail/"
    err := database.DB.QueryRow(`
      SELECT COUNT(*) 
      FROM logs 
      WHERE timestamp >= $1
        AND endpoint LIKE '/detail/%'
    `, startOfDay).Scan(&count)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"linkClicks": count})
}
