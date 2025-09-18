package routes

import "github.com/gin-gonic/gin"

func StartRoutes(server *gin.Engine) {
	server.POST("/bookstore/books", CreateBook)
	server.GET("/bookstore/books", GetBooks)
	server.GET("/bookstore/books/findByTitle/:title", GetBookByTitle)
	server.GET("/bookstore/books/findByAuthor/:author", GetBooksByAuthor)
	server.POST("/bookstore/users/signup", signup)
	server.POST("/bookstore/users/login", login)
}