package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	r "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
)

type loanRepository struct {
	db *sql.DB
}

// NewLoanRepository creates a new instance of loanRepository
func NewLoanRepository(db *sql.DB) r.LoanRepository {
	return &loanRepository{
		db: db,
	}
}

// CheckNotReturned verify if a loan not returned exists in the database
func (r *loanRepository) CheckNotReturned(ctx context.Context, userID int, bookID int) (bool, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT * FROM loans WHERE is_returned = false AND user_id = $1 AND book_id = $2")
	if err != nil {
		return false, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}
	defer stmt.Close()

	l := &entity.Loan{}

	row := stmt.QueryRowContext(ctx, userID, bookID)

	err = row.Scan(&l.ID, &l.UserID, &l.BookID, &l.Is_returned, &l.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// all books are returned
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", ErrScanData, err)
	}

	return true, nil
}

// Search searches all books a user borrowed
func (r *loanRepository) Search(ctx context.Context, userID int) ([]*entity.Loan, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT * FROM loans WHERE user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrExecuteQuery, err)
	}

	var loans []*entity.Loan
	for rows.Next() {
		var l entity.Loan

		err = rows.Scan(&l.ID, &l.UserID, &l.BookID, &l.Is_returned, &l.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrScanData, err)
		}

		loans = append(loans, &l)
	}

	if len(loans) == 0 {
		return nil, ErrLoanNotFound
	}

	return loans, nil
}

// BorrowTransaction borrows a book updating the book amount and creating a new loan
func (r *loanRepository) BorrowTransaction(ctx context.Context, u *entity.User, b *entity.Book) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrBeginTransaction, err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE books SET amount = $1 WHERE id = $2", b.Amount, b.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%s: %w", ErrRollback, rbErr)
		}
		return fmt.Errorf("%s: %w", ErrExecuteStatement, err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO loans (user_id, book_id) VALUES ($1, $2)", u.ID, b.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%s: %w", ErrRollback, rbErr)
		}
		return fmt.Errorf("%s: %w", ErrExecuteStatement, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", ErrCommit, err)
	}

	return nil
}

// ReturnTransaction returns a book updating the book amount and updating the existing loan
func (r *loanRepository) ReturnTransaction(ctx context.Context, u *entity.User, b *entity.Book) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrBeginTransaction, err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE books SET amount = $1 WHERE id = $2", b.Amount, b.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%s: %w", ErrRollback, rbErr)
		}
		return fmt.Errorf("%s: %w", ErrExecuteStatement, err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE loans SET is_returned = $1 WHERE user_id = $2 AND book_id = $3", true, u.ID, b.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%s: %w", ErrRollback, rbErr)
		}
		return fmt.Errorf("%s: %w", ErrExecuteStatement, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", ErrCommit, err)
	}

	return nil
}
