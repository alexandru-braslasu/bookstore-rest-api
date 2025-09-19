package routes

import (
	"net/http"
	"net/mail"

	"bookstore.com/booking/googleauth"
	"bookstore.com/booking/models"
	"bookstore.com/booking/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBind(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request."})
		return
	}

	_, err = mail.ParseAddress(user.Email)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid address format!"})
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

func loginWithGoogle(context *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	googleauth.RegisterProvider()
	query := context.Request.URL.Query()
    query.Add("provider", "google")
    context.Request.URL.RawQuery = query.Encode()
    gothic.BeginAuthHandler(context.Writer, context.Request)
}

func AuthResponse(context *gin.Context) {
    query := context.Request.URL.Query()
    query.Add("provider", "google")
    context.Request.URL.RawQuery = query.Encode()

    user, err := gothic.CompleteUserAuth(context.Writer, context.Request)
    if err != nil {
        context.AbortWithError(http.StatusInternalServerError, err)
        return
    }

	exists, err := models.VerifyExists(user.Email)
	
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
        return
	}

	if exists == 0 {
		var userInDB models.User
		userInDB.Email = user.Email
		userInDB.Name = user.Name
		userInDB.Password = ""
		userInDB.IsAdmin = false
		userInDB.AddUserToTable() 
	}

	context.JSON(http.StatusOK, user)
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