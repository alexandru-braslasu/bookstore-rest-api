package routes

import (
	"net/http"

	"bookstore.com/booking/models"
	"github.com/gin-gonic/gin"
)

func CreateBook(context *gin.Context) {
	isAdmin := context.GetBool("isAdmin")

	if !isAdmin {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized operation! :("})
		return
	}

	var book models.Book
	err := context.ShouldBind(&book)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request"})
		return
	}

	// models.Books = append(models.Books, book)

	err = book.AddBookInTable()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create book. Try again later!"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Book created successfully!"})
}

func GetBooks(context *gin.Context) {
	// books := models.Books

	books, err := models.GetAllBooks()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get books."})
		return
	}

	context.JSON(http.StatusOK, books)
}

func GetBookByTitle(context *gin.Context) {
	title := context.Param("title")

	// var bookFound models.Book
	// found := false
	// for _, book := range models.Books {
	// 	if strings.Compare(book.Title, title) == 0 {
	// 		bookFound = book
	// 		found = true
	// 	}
	// }

	// if !found {
	// 	context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't find the book. Try again later!"})
	// 	return
	// }

	bookFound, err := models.GetBookFromTable(title)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find book. Try again later!"})
		return
	}

	context.JSON(http.StatusOK, bookFound)
}

func GetBooksByAuthor(context *gin.Context) {
	author := context.Param("author")

	// var booksFound []models.Book
	// for _, book := range models.Books {
	// 	if strings.Contains(book.Author, author) {
	// 		booksFound = append(booksFound, book)
	// 	}
	// }

	booksFound, err := models.GetBooksFromTable(author)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find books with that author."})
		return
	}

	context.JSON(http.StatusOK, booksFound)
}