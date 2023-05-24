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

func TestCheckNotReturned(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	loan := &entity.Loan{
		ID:          1,
		BookID:      1,
		UserID:      1,
		Is_returned: false,
		CreatedAt:   time.Now(),
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "is_returned", "amount", "created_at"}).
			AddRow(loan.ID, loan.UserID, loan.BookID, loan.Is_returned, loan.CreatedAt)

		mock.ExpectPrepare("SELECT \\* FROM loans WHERE is_returned = false AND user_id = \\$1 AND book_id = \\$2").
			ExpectQuery().
			WithArgs(loan.UserID, loan.BookID).
			WillReturnRows(rows)

		exists, err := repo.CheckNotReturned(context.Background(), loan.UserID, loan.BookID)
		assert.NoError(t, err)
		assert.True(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM loans WHERE is_returned = false AND user_id = \\$1 AND book_id = \\$2").
			WillReturnError(sql.ErrConnDone)

		exists, err := repo.CheckNotReturned(context.Background(), loan.UserID, loan.BookID)
		assert.Error(t, err)
		assert.False(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM loans WHERE is_returned = false AND user_id = \\$1 AND book_id = \\$2").
			ExpectQuery().
			WillReturnError(sql.ErrConnDone)

		exists, err := repo.CheckNotReturned(context.Background(), loan.UserID, loan.BookID)
		assert.Error(t, err)
		assert.False(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM loans WHERE is_returned = false AND user_id = \\$1 AND book_id = \\$2").
			ExpectQuery().
			WithArgs(10, 10).
			WillReturnError(sql.ErrNoRows)

		exists, err := repo.CheckNotReturned(context.Background(), 10, 10)
		assert.NoError(t, err)
		assert.False(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSearch(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	loans := []*entity.Loan{
		{
			ID:          1,
			UserID:      1,
			BookID:      1,
			Is_returned: false,
			CreatedAt:   time.Now(),
		},
		{
			ID:          2,
			UserID:      1,
			BookID:      2,
			Is_returned: true,
			CreatedAt:   time.Now(),
		},
		{
			ID:          3,
			UserID:      1,
			BookID:      10,
			Is_returned: false,
			CreatedAt:   time.Now(),
		},
	}
	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "is_returned", "amount", "created_at"})
		for _, loan := range loans {
			rows = rows.AddRow(loan.ID, loan.UserID, loan.BookID, loan.Is_returned, loan.CreatedAt)
		}

		mock.ExpectPrepare("SELECT \\* FROM loans WHERE user_id =").
			ExpectQuery().
			WithArgs(1).
			WillReturnRows(rows)

		gotLoans, err := repo.Search(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, loans)

		for i := range gotLoans {
			assert.Equal(t, loans[i], gotLoans[i])
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM loans WHERE user_id =").
			WillReturnError(sql.ErrConnDone)

		gotLoans, err := repo.Search(context.Background(), 1)
		assert.Error(t, err)
		assert.Empty(t, gotLoans)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM loans WHERE user_id =").
			ExpectQuery().
			WillReturnError(sql.ErrConnDone)

		gotLoans, err := repo.Search(context.Background(), 1)
		assert.Error(t, err)
		assert.Empty(t, gotLoans)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM loans WHERE user_id =").
			ExpectQuery().
			WithArgs(10).
			WillReturnRows(&sqlmock.Rows{})

		gotLoans, err := repo.Search(context.Background(), 10)
		assert.Error(t, err)
		assert.Empty(t, gotLoans)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBorrowTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	user := &entity.User{
		ID: 1,
	}

	book := &entity.Book{
		ID:     1,
		Amount: 5,
	}

	t.Run("OK", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE books SET amount = \\$1 WHERE id = \\$2").
			WithArgs(book.Amount, book.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO loans \\(user_id, book_id\\) VALUES \\(\\$1, \\$2\\)").
			WithArgs(user.ID, book.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.BorrowTransaction(context.Background(), user, book)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Exec Update Books Failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE books SET amount = \\$1 WHERE id = \\$2").
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.BorrowTransaction(context.Background(), user, book)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Exec Insert Loans Failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE books SET amount = \\$1 WHERE id = \\$2").
			WithArgs(book.Amount, book.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO loans \\(user_id, book_id\\) VALUES \\(\\$1, \\$2\\)").
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		err := repo.BorrowTransaction(context.Background(), user, book)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReturnTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	user := &entity.User{
		ID: 1,
	}

	book := &entity.Book{
		ID:     1,
		Amount: 5,
	}

	t.Run("OK", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE books SET amount = \\$1 WHERE id = \\$2").
			WithArgs(book.Amount, book.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE loans SET is_returned = \\$1 WHERE user_id = \\$2 AND book_id = \\$3").
			WithArgs(true, user.ID, book.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.ReturnTransaction(context.Background(), user, book)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Exec Update Books Failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE books SET amount = \\$1 WHERE id = \\$2").
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.ReturnTransaction(context.Background(), user, book)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Exec Update Loans Failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE books SET amount = \\$1 WHERE id = \\$2").
			WithArgs(book.Amount, book.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE loans SET is_returned = \\$1 WHERE user_id = \\$2 AND book_id = \\$3").
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.ReturnTransaction(context.Background(), user, book)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
