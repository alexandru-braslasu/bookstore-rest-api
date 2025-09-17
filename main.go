package main

import (
	"bookstore.com/booking/db"
	"bookstore.com/booking/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	routes.StartRoutes(server)

	server.Run(":8080")
}