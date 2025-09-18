package models

import (
	"errors"
	"fmt"

	"bookstore.com/booking/db"
	"bookstore.com/booking/utils"
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

func (user *User) GetIsAdminByEmail() error {
	query := `
	SELECT isAdmin
	FROM users
	WHERE email = ?
	`
	row := db.DB.QueryRow(query, user.Email)
	var isAdmin bool
	err := row.Scan(&isAdmin)

	if err != nil {
		return err
	}

	user.IsAdmin = isAdmin
	return nil
}

func GetBooksSavedByUser(email string) ([]Book, error) {
	query := `
	SELECT title, author, description, nrSamples
	FROM saved_books
	WHERE email = ?
	`
	rows, err := db.DB.Query(query, email)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fmt.Println("here")

	var books []Book
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.Title, &book.Author, &book.Description, &book.NrSamples)

		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}
	return books, nil
}