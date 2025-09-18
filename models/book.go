package models

import "bookstore.com/booking/db"

type Book struct {
	Title       string `binding:"required"`
	Author      string `binding:"required"`
	Description string
	NrSamples   int `binding:"required"`
}

// var Books []Book

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

func (book *Book) UpdateBook() error {
	query := `
	UPDATE books
	SET nrSamples = nrSamples - ?
	WHERE title = ? AND author = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(book.NrSamples, book.Title, book.Author)

	return err
}

func (book *Book) GetDescription() error {
	query := `
	SELECT description
	FROM books
	WHERE title = ? AND author = ?
	`
	row := db.DB.QueryRow(query, book.Title, book.Author)
	var description string
	err := row.Scan(&description)
	
	if err != nil {
		return err
	}

	book.Description = description
	return nil
}

func (book *Book) GetNrSamples() (int, error) {
	query := `
	SELECT nrSamples
	FROM books
	WHERE title = ? AND author = ?
	`
	row := db.DB.QueryRow(query, book.Title, book.Author)
	var nrSamples int
	err := row.Scan(&nrSamples)
	
	if err != nil {
		return 0, err
	}

	return nrSamples, nil
}

func (book Book) EliminateBookFromStore() error {
	query := `
	DELETE FROM books
	WHERE title = ? AND author = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(book.Title, book.Author)

	return err
}

func (book Book) AddToSaved(email string) error {
	exists, err := book.verifyIfExists()
	
	if err != nil {
		return err
	}

	if exists == 0 {
		query := `
		INSERT INTO saved_books(title, author, description, email, nrSamples)
		VALUES (?, ?, ?, ?, ?)
		`
		stmt, err := db.DB.Prepare(query)

		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(book.Title, book.Author, book.Description, email, book.NrSamples)

		if err != nil {
			return err
		}
		return nil
	}
	query := `
	UPDATE saved_books
	SET nrSamples = nrSamples + ?
	WHERE title = ? AND author = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(book.NrSamples, book.Title, book.Author)

	if err != nil {
		return err
	}

	return nil
}

func (book Book) verifyIfExists() (int, error) {
	query := `
	SELECT COUNT(*)
	FROM saved_books
	WHERE title = ? AND author = ?
	`
	row := db.DB.QueryRow(query, book.Title, book.Author)
	var exists int
	err := row.Scan(&exists)

	if err != nil {
		return 0, err
	}

	return exists, nil
}