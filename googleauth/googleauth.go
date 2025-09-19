package googleauth

import (
	"os"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func RegisterProvider() {
    goth.UseProviders(
        google.New(
            os.Getenv("CLIENT_ID"),
            os.Getenv("CLIENT_SECRET"),
            "http://localhost:8080/bookstore/users/auth",
        ),
    )
}