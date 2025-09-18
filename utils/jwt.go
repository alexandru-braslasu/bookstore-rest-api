package utils

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func GenerateToken(email string, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"isAdmin": isAdmin,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(viperGetSecretKey()))
}

func VerifyToken(token string) (string, bool, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(viperGetSecretKey()), nil
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

func viperGetSecretKey() string {
  viper.SetConfigFile(".env")

  err := viper.ReadInConfig()

  if err != nil {
    log.Fatalf("Error while reading config file %s", err)
  }

  value, ok := viper.Get("SECRET_KEY").(string)

  if !ok {
    log.Fatalf("Invalid type assertion")
  }

  return value
}