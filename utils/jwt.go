package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateToken(email string, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"isAdmin": isAdmin,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func VerifyToken(token string) (string, bool, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return "", false, errors.New("Could not parse token.")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return "", false, errors.New("Invalid token!")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", false, errors.New("Invalid token claims.")
	}

	email := claims["email"].(string)
	isAdmin := claims["isAdmin"].(bool)

	return email, isAdmin, nil
}

type TokenInfo struct {
    Issuer         string `json:"iss"`
    Audience       string `json:"aud"`
    ExpirationTime string `json:"exp"`
    IssuedTime     string `json:"iat"`
    Email          string `json:"email"`
    EmailVerified  string `json:"email_verified"`
}

func VerifyGoogleToken(id_token string) (string, bool, error) {
	const googleTokenURL = "https://oauth2.googleapis.com/tokeninfo?id_token=%s"
	err := godotenv.Load()
	if err != nil {
		return "", false, err
	}
    resp, err := http.Get(fmt.Sprintf(googleTokenURL, id_token))
    if err != nil {
        return "", false, err
    }
    defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {
            log.Println(err)
        }
    }(resp.Body)

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", false, err
    }

    var tokenInfo TokenInfo
    err = json.Unmarshal(body, &tokenInfo)
    if err != nil {
        return "", false, err
    }
    if !strings.Contains(tokenInfo.Audience, os.Getenv("CLIENT_ID")) {
        return "", false, errors.New("token err")
    }

    return tokenInfo.Email, false, nil
}