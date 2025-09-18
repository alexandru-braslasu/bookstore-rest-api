package routes

import (
	"net/http"

	"bookstore.com/booking/models"
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

	context.JSON(http.StatusOK, gin.H{"message": "Logged in successfully!"})
}