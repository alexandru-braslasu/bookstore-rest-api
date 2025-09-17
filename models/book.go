package models

import "bookstore.com/booking/db"

type Book struct {
	Title       string `binding:"required"`
	Author      string `binding:"required"`
	Description string
	NrSamples   int `binding:"required"`
}

var Books []Book

func (book Book) AddBookInTable() error {
	query := `
	INSERT INTO books(title, author, description, nrSamples)
	VALUES (?, ?, ?, ?)
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(book.Title, book.Author, book.Description, book.NrSamples)
	return err
}

func GetAllBooks() ([]Book, error) {
	query := `
	SELECT * FROM books
	`
	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Title, &book.Author, &book.Description, &book.NrSamples)

		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}
	return books, nil
}

func GetBookFromTable(title string) (*Book, error) {
	query := `
	SELECT *
	FROM books
	WHERE title = ?
	`
	row := db.DB.QueryRow(query, title)
	var book Book
	err := row.Scan(&book.Title, &book.Author, &book.Description, &book.NrSamples)
	if err != nil {
		return nil, err
	}
	return &book, err
}

func GetBooksFromTable(author string) ([]Book, error) {
	query := `
	SELECT *
	FROM books
	WHERE author LIKE ?
	`
	rows, err := db.DB.Query(query, "%" + author + "%")

	if err != nil {
		return nil, err
	}

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Title, &book.Author, &book.Description, &book.NrSamples)

		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}
	return books, nil
}