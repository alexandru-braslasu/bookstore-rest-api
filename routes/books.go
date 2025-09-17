package routes

import (
	"net/http"
	"strings"

	"bookstore.com/booking/models"
	"github.com/gin-gonic/gin"
)

func CreateBook(context *gin.Context) {
	var book models.Book
	err := context.ShouldBind(&book)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request"})
		return
	}

	models.Books = append(models.Books, book)
	context.JSON(http.StatusCreated, gin.H{"message": "Book created successfully!"})
}

func GetBooks(context *gin.Context) {
	books := models.Books
	context.JSON(http.StatusOK, books)
}

func GetBookByTitle(context *gin.Context) {
	title := context.Param("title")

	var bookFound models.Book
	found := false
	for _, book := range models.Books {
		if strings.Compare(book.Title, title) == 0 {
			bookFound = book
			found = true
		}
	}

	if !found {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't find the book. Try again later!"})
		return
	}

	context.JSON(http.StatusOK, bookFound)
}

func GetBooksByAuthor(context *gin.Context) {
	author := context.Param("author")

	var booksFound []models.Book
	for _, book := range models.Books {
		if strings.Contains(book.Author, author) {
			booksFound = append(booksFound, book)
		}
	}

	context.JSON(http.StatusOK, booksFound)
}