package main

import (
	"project/database"
	"project/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type","apikey","X-Email"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	  }))
	  

	r.POST("/register", handlers.RegisterClient)
	r.GET("/client",handlers.GetClient)
	auth := r.Group("/")
	auth.Use(handlers.AuthMiddleware(), handlers.LogMiddleware())
	{
		auth.GET("/data", handlers.GetData)
		auth.GET("/data/:id", handlers.GetDataByID)
		auth.GET("/search", handlers.GetDataByQuery)
		auth.GET("/log",handlers.GetLog)
		auth.GET("log/:id",handlers.GetLogByID)
		auth.GET("/stats/clicks-search",handlers.GetClicksSearch)
		auth.POST("/log/click",handlers.PostClickLog)
		auth.GET("/stats/link-clicks",handlers.GetLinkClicks)
	}

	r.Run(":8080")
}
