package bookRepository

import (
	"books-app/models"
	"database/sql"
	"log"
)

type BookRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (b BookRepository) GetBooks(db *sql.DB, book models.Book, books []models.Book) []models.Book {
	rows, err := db.Query("select * from books where user_id=$1", book.UserId)
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.UserId)
		logFatal(err)

		books = append(books, book)
	}

	return books
}

func (b BookRepository) GetBook(db *sql.DB, book models.Book, id int) models.Book {
	rows := db.QueryRow("select * from books where id=$1 and user_id=$2", id, book.UserId)

	err := rows.Scan(&book.ID, &book.UserId, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	log.Println(book)

	return book
}

func (b BookRepository) AddBook(db *sql.DB, book models.Book) int {
	err := db.QueryRow("insert into books (title, author, year, user_id) values($1, $2, $3, $4) RETURNING id;",
		book.Title, book.Author, book.Year, book.UserId).Scan(&book.ID)

	logFatal(err)

	return book.ID
}

func (b BookRepository) UpdateBook(db *sql.DB, book models.Book) int64 {
	result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id",
		&book.Title, &book.Author, &book.Year, &book.ID)

	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	return rowsUpdated
}

func (b BookRepository) RemoveBook(db *sql.DB, id int) int64 {
	result, err := db.Exec("delete from books where id = $1", id)
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	return rowsDeleted
}
