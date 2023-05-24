package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
)

func TestGetBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepository(db)

	book := &entity.Book{
		ID:        1,
		Title:     "Let's Go Further!",
		Author:    "Alex Edwards",
		Amount:    5,
		UpdatedAt: time.Time{},
		CreatedAt: time.Now(),
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author", "amount", "updated_at", "created_at"}).
			AddRow(book.ID, book.Title, book.Author, book.Amount, book.UpdatedAt, book.CreatedAt)

		mock.ExpectPrepare("SELECT \\* FROM books WHERE id = ").
			ExpectQuery().
			WithArgs(book.ID).
			WillReturnRows(rows)

		gotBook, err := repo.Get(context.Background(), book.ID)
		assert.NoError(t, err)
		assert.Equal(t, book, gotBook)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM books WHERE id = ").
			WillReturnError(sql.ErrConnDone)

		gotBook, err := repo.Get(context.Background(), book.ID)
		assert.Error(t, err)
		assert.Empty(t, gotBook)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM books WHERE id = ").
			ExpectQuery().
			WithArgs(book.ID).
			WillReturnError(sql.ErrConnDone)

		gotBook, err := repo.Get(context.Background(), book.ID)
		assert.Error(t, err)
		assert.Empty(t, gotBook)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM books WHERE id = ").
			ExpectQuery().
			WithArgs(book.ID).
			WillReturnError(sql.ErrNoRows)

		gotBook, err := repo.Get(context.Background(), book.ID)
		assert.Error(t, err)
		assert.Empty(t, gotBook)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestListBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepository(db)

	books := []*entity.Book{
		{
			ID:        1,
			Title:     "Let's Go Further!",
			Author:    "Alex Edwards",
			Amount:    5,
			UpdatedAt: time.Time{},
			CreatedAt: time.Now(),
		}, {
			ID:        2,
			Title:     "Let's Go!",
			Author:    "Alex Edwards",
			Amount:    10,
			UpdatedAt: time.Time{},
			CreatedAt: time.Now(),
		},
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author", "amount", "updated_at", "created_at"})
		for _, book := range books {
			rows = rows.AddRow(book.ID, book.Title, book.Author, book.Amount, book.UpdatedAt, book.CreatedAt)
		}

		mock.ExpectPrepare("SELECT \\* FROM books").
			ExpectQuery().
			WillReturnRows(rows)

		gotBooks, err := repo.List(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, len(books), len(gotBooks))

		for i := range gotBooks {
			assert.Equal(t, books[i], gotBooks[i])
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM books").
			WillReturnError(sql.ErrConnDone)

		gotBook, err := repo.List(context.Background())
		assert.Error(t, err)
		assert.Empty(t, gotBook)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM books").
			ExpectQuery().
			WillReturnError(sql.ErrConnDone)

		gotBook, err := repo.List(context.Background())
		assert.Error(t, err)
		assert.Empty(t, gotBook)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Empty book", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM books").
			ExpectQuery().
			WillReturnRows(&sqlmock.Rows{})

		gotBooks, err := repo.List(context.Background())
		assert.Error(t, err)
		assert.Nil(t, gotBooks)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSearchBooks(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepository(db)

	books := []*entity.Book{
		{
			ID:        1,
			Title:     "Let's Go Further!",
			Author:    "Alex Edwards",
			Amount:    5,
			UpdatedAt: time.Time{},
			CreatedAt: time.Now(),
		}, {
			ID:        2,
			Title:     "Let's Go!",
			Author:    "Alex Edwards",
			Amount:    10,
			UpdatedAt: time.Time{},
			CreatedAt: time.Now(),
		},
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author", "amount", "updated_at", "created_at"})
		for _, book := range books {
			rows = rows.AddRow(book.ID, book.Title, book.Author, book.Amount, book.UpdatedAt, book.CreatedAt)
		}

		mock.ExpectPrepare("SELECT \\* FROM books WHERE LOWER\\(title\\) LIKE LOWER\\(\\$1\\)").
			ExpectQuery().
			WithArgs("%Let's Go%").
			WillReturnRows(rows)

		gotBooks, err := repo.Search(context.Background(), "Let's Go")

		assert.NoError(t, err)
		assert.Equal(t, len(books), len(gotBooks))

		for i := range gotBooks {
			assert.Equal(t, books[i], gotBooks[i])
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM books WHERE LOWER\\(title\\) LIKE LOWER\\(\\$1\\)").
			WillReturnError(sql.ErrConnDone)

		gotBooks, err := repo.Search(context.Background(), "Let's Go")

		assert.Error(t, err)
		assert.Empty(t, gotBooks)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM books WHERE LOWER\\(title\\) LIKE LOWER\\(\\$1\\)").
			ExpectQuery().
			WillReturnError(sql.ErrConnDone)

		gotBooks, err := repo.Search(context.Background(), "Let's Go")

		assert.Error(t, err)
		assert.Empty(t, gotBooks)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Empty Book", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM books WHERE LOWER\\(title\\) LIKE LOWER\\(\\$1\\)").
			ExpectQuery().
			WillReturnRows(&sqlmock.Rows{})

		gotBooks, err := repo.Search(context.Background(), "Let's Go")

		assert.Error(t, err)
		assert.Empty(t, gotBooks)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM books WHERE LOWER\\(title\\) LIKE LOWER\\(\\$1\\)").
			ExpectQuery().
			WithArgs("%Let's Go%").
			WillReturnRows(&sqlmock.Rows{})

		gotBooks, err := repo.Search(context.Background(), "Let's Go")

		assert.Error(t, err)
		assert.Empty(t, gotBooks)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepository(db)

	book := &entity.Book{
		ID:        1,
		Title:     "Let's Go Further!",
		Author:    "Alex Edwards",
		Amount:    5,
		UpdatedAt: time.Time{},
		CreatedAt: time.Now(),
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(1)

		mock.ExpectPrepare("INSERT INTO books").
			ExpectQuery().
			WithArgs(book.Title, book.Author, book.Amount).
			WillReturnRows(rows)

		id, err := repo.Create(context.Background(), book)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO books").
			WillReturnError(sql.ErrConnDone)

		id, err := repo.Create(context.Background(), book)
		assert.Error(t, err)
		assert.Empty(t, id)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO books").
			ExpectQuery().
			WillReturnError(sql.ErrConnDone)

		id, err := repo.Create(context.Background(), book)
		assert.Error(t, err)
		assert.Empty(t, id)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepository(db)

	book := &entity.Book{
		ID:        1,
		Title:     "Let's Go Further!",
		Author:    "Alex Edwards",
		Amount:    5,
		UpdatedAt: time.Time{},
		CreatedAt: time.Now(),
	}

	t.Run("OK", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE books").
			ExpectExec().
			WithArgs(book.Title, book.Author, book.Amount, book.ID).
			WillReturnResult(sqlmock.NewResult(int64(book.ID), 1))

		err := repo.Update(context.Background(), book)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE books").
			WillReturnError(sql.ErrConnDone)

		err := repo.Update(context.Background(), book)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE books").
			ExpectExec().
			WillReturnError(sql.ErrConnDone)

		err := repo.Update(context.Background(), book)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE books").
			ExpectExec().
			WithArgs(book.Title, book.Author, book.Amount, book.ID).
			WillReturnResult(sqlmock.NewResult(int64(book.ID), 0))

		err := repo.Update(context.Background(), book)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepository(db)

	book := &entity.Book{
		ID:        1,
		Title:     "Let's Go Further!",
		Author:    "Alex Edwards",
		Amount:    5,
		UpdatedAt: time.Time{},
		CreatedAt: time.Now(),
	}

	t.Run("OK", func(t *testing.T) {
		mock.ExpectPrepare("DELETE FROM books WHERE id =").
			ExpectExec().
			WithArgs(book.ID).
			WillReturnResult(sqlmock.NewResult(int64(book.ID), 1))

		err := repo.Delete(context.Background(), book.ID)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("DELETE FROM books WHERE id =").
			WillReturnError(sql.ErrConnDone)

		err := repo.Delete(context.Background(), book.ID)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Exec Failed", func(t *testing.T) {
		mock.ExpectPrepare("DELETE FROM books WHERE id =").
			ExpectExec().
			WillReturnError(sql.ErrConnDone)

		err := repo.Delete(context.Background(), book.ID)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("DELETE FROM books WHERE id =").
			ExpectExec().
			WithArgs(book.ID).
			WillReturnResult(sqlmock.NewResult(int64(book.ID), 0))

		err := repo.Delete(context.Background(), book.ID)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
