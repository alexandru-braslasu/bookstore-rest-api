package models

import (
	"errors"
	"log"

	"bookstore.com/booking/db"
	"bookstore.com/booking/utils"
	"github.com/spf13/viper"
)

type User struct {
	Name     string
	Email    string `binding:"required"`
	Password string `binding:"required"`
	IsAdmin  bool
}

func (user User) AddUserToTable() error {
	query := `
	INSERT INTO users(name, email, password, isAdmin)
	VALUES (?, ?, ?, ?)
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.Name, user.Email, hashedPassword, user.IsAdmin)
	return err
}

func (user User) ValidateUser() error {
	query := `
	SELECT password
	FROM users
	WHERE email = ?
	`
	row := db.DB.QueryRow(query, user.Email)
	var retrievedPassword string
	err := row.Scan(&retrievedPassword)
	
	if err != nil{
		return errors.New("Invalid credentials!")
	}

	ok := utils.CheckPasswordHash(user.Password, retrievedPassword)

	if !ok {
		return errors.New("Invalid credentials!")
	}

	return nil
}

func ViperGetSecretKey() string {
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