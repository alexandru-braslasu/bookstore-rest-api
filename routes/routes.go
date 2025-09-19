package routes

import (
	"bookstore.com/booking/middlewares"
	"github.com/gin-gonic/gin"
)

func StartRoutes(server *gin.Engine) {
	server.GET("/bookstore/books", GetBooks)
	server.GET("/bookstore/books/findByTitle/:title", GetBookByTitle)
	server.GET("/bookstore/books/findByAuthor/:author", GetBooksByAuthor)
	server.POST("/bookstore/users/signup", signup)
	server.POST("/bookstore/users/login", login)
	server.GET("/bookstore/users/loginWithGoogle", loginWithGoogle)
	server.GET("/bookstore/users/auth", AuthResponse)

	authenticated := server.Group("/bookstore")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/books", CreateBook)
	authenticated.PUT("/users/saveBook", saveBook)
	authenticated.GET("/users/savedBooks", GetBooksSaved)
}