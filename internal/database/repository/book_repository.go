package repository

import (
	"database/sql"
	"errors"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	r "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
)

type bookRepository struct {
	db *sql.DB
}

// NewBookRepository creates a new instance of BookRepository
func NewBookRepository(db *sql.DB) r.BookRepository {
	return &bookRepository{
		db: db,
	}
}

// Get gets book info by id
func (r *bookRepository) Get(id int) (*entity.Book, error) {
	stmt, err := r.db.Prepare("SELECT * FROM books WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	b := &entity.Book{}

	row := stmt.QueryRow(id)

	var updatedAt sql.NullTime
	err = row.Scan(&b.ID, &b.Title, &b.Author, &b.Amount, &updatedAt, &b.CreatedAt)
	if err != nil {
		return nil, err
	}

	// check if updated_at is NULL before scanning it
	if updatedAt.Valid {
		b.UpdatedAt = updatedAt.Time
	}

	return b, nil
}

// List list all books in the database
func (r *bookRepository) List() ([]*entity.Book, error) {
	stmt, err := r.db.Prepare("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var books []*entity.Book
	for rows.Next() {
		var b entity.Book
		var updatedAt sql.NullTime

		err = rows.Scan(&b.ID, &b.Title, &b.Author, &b.Amount, &updatedAt, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		// check if updated_at is NULL before scanning it
		if updatedAt.Valid {
			b.UpdatedAt = updatedAt.Time
		}

		books = append(books, &b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// Search searches books matching the sent query
func (r *bookRepository) Search(query string) ([]*entity.Book, error) {
	stmt, err := r.db.Prepare("SELECT * FROM books WHERE LOWER(title) LIKE LOWER($1)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + query + "%")
	if err != nil {
		return nil, err
	}

	var books []*entity.Book
	for rows.Next() {
		var b entity.Book
		var updatedAt sql.NullTime

		err = rows.Scan(&b.ID, &b.Title, &b.Author, &b.Amount, &updatedAt, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		// check if updated_at is NULL before scanning it
		if updatedAt.Valid {
			b.UpdatedAt = updatedAt.Time
		}

		books = append(books, &b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// Create creates a new book
func (r *bookRepository) Create(b *entity.Book) (int, error) {
	stmt, err := r.db.Prepare("INSERT INTO books (title, author, amount) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(b.Title, b.Author, b.Amount).Scan(&b.ID)
	if err != nil {
		return 0, err
	}

	return b.ID, nil
}

// Update updates a book
func (r *bookRepository) Update(b *entity.Book) error {
	stmt, err := r.db.Prepare("UPDATE books SET title = $1, author = $2, amount = $3, updated_at = $4 WHERE id = $5")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(b.Title, b.Author, b.Amount, b.UpdatedAt, b.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}

// Delete deletes a book by id
func (r *bookRepository) Delete(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM books WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}
