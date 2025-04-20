package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

func GetData(c *gin.Context) {
    rows, err := database.DB.Query("SELECT id, content FROM data")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch data"})
        return
    }
    defer rows.Close()

    // Initialize เป็น slice ว่างแทนที่จะเป็น nil
    results := make([]models.Data, 0)
    for rows.Next() {
        var data models.Data
        if err := rows.Scan(&data.ID, &data.Content); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Data scan error"})
            return
        }
        results = append(results, data)
    }

    c.JSON(http.StatusOK, results)
}

func GetDataByID(c *gin.Context) {
    id := c.Param("id")

    // ใช้ $1 สำหรับ PostgreSQL
    row := database.DB.QueryRow(
        "SELECT id, content FROM data WHERE id = $1",
        id,
    )

    var data models.Data
    err := row.Scan(&data.ID, &data.Content)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
        } else {
            // แสดง error รายละเอียดเข้า log ก่อนก็ได้
            log.Printf("GetDataByID scan error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Data scan error"})
        }
        return
    }

    c.JSON(http.StatusOK, data)
}


func GetDataByQuery(c *gin.Context) {
    // อ่าน q จาก URL ?q=...
    keyword := c.Query("keyword")

    var (
        rows *sql.Rows
        err  error
    )

    if keyword != "" {
        // ถ้าใช้ PostgreSQL ให้ใช้ ILIKE, ถ้า MySQL ใช้ LIKE
        search := "%" + keyword + "%"
        rows, err = database.DB.Query(
            "SELECT id, content FROM data WHERE content ILIKE $1",
            search,
        )
    } else {
        rows, err = database.DB.Query("SELECT id, content FROM data")
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch data"})
        return
    }
    defer rows.Close()

    results := make([]models.Data, 0)
    for rows.Next() {
        var data models.Data
        if err := rows.Scan(&data.ID, &data.Content); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Data scan error"})
            return
        }
        results = append(results, data)
    }
    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Row iteration error"})
        return
    }

    c.JSON(http.StatusOK, results)
}