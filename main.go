package main

import (
	"bookstore.com/booking/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	routes.StartRoutes(server)

	server.Run(":8080")
}