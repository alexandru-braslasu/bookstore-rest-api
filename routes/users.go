package routes

import (
	"net/http"

	"bookstore.com/booking/models"
	"bookstore.com/booking/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBind(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request."})
		return
	}

	err = user.AddUserToTable()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add user. Try again later!"})
		return
	}
	
	context.JSON(http.StatusCreated, gin.H{"message": "User added successfully!"})
}

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBind(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request."})
		return
	}

	err = user.ValidateUser()

	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = user.GetIsAdminByEmail()

	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	}

	token, err := utils.GenerateToken(user.Email, user.IsAdmin)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": token})
}

func saveBook(context *gin.Context) {
	var book models.Book
	err := context.ShouldBind(&book)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request"})
		return
	}

	nrSamples, err := book.GetNrSamples()

	if err != nil || nrSamples < book.NrSamples {
		context.JSON(http.StatusBadRequest, gin.H{"message": "There are not enough samples. Try again later!"})
		return
	}

	err = book.GetDescription()

	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = book.UpdateBook()

	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	nrSamplesRemained, err := book.GetNrSamples()

	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if nrSamplesRemained == 0 {
		err = book.EliminateBookFromStore()
		if err != nil {
			context.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	err = book.AddToSaved(context.GetString("email"))

	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Book added successfully."})
}

func GetBooksSaved(context *gin.Context) {
	email := context.GetString("email")
	books, err := models.GetBooksSavedByUser(email)
	
	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	context.JSON(http.StatusOK, books)
}