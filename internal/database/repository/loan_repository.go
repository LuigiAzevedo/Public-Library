package repository

import (
	"database/sql"

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
func (r *loanRepository) CheckNotReturned(userID int, bookID int) (bool, error) {
	stmt, err := r.db.Prepare("SELECT * FROM loans WHERE is_returned = false AND user_id = $1 AND book_id = $2")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	l := &entity.Loan{}

	row := stmt.QueryRow(userID, bookID)
	err = row.Scan(&l.ID, &l.UserID, &l.BookID, &l.Is_returned, &l.CreatedAt)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Search searches all books a user borrowed
func (r *loanRepository) Search(userID int) ([]*entity.Loan, error) {
	stmt, err := r.db.Prepare("SELECT * FROM loans WHERE user_id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	var loans []*entity.Loan
	for rows.Next() {
		var l entity.Loan

		err = rows.Scan(&l.ID, &l.UserID, &l.BookID, &l.Is_returned, &l.CreatedAt)
		if err != nil {
			return nil, err
		}

		loans = append(loans, &l)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}

// BorrowTransaction borrows a book updating the book amount and creating a new loan
func (r *loanRepository) BorrowTransaction(u *entity.User, b *entity.Book) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE books SET amount = $1 WHERE id = $2", b.Amount, b.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	_, err = tx.Exec("INSERT INTO loans (user_id, book_id) VALUES ($1, $2)", u.ID, b.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ReturnTransaction returns a book updating the book amount and updating the existing loan
func (r *loanRepository) ReturnTransaction(u *entity.User, b *entity.Book) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE books SET amount = $1 WHERE id = $2", b.Amount, b.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	_, err = tx.Exec("UPDATE loans SET is_returned = $1 WHERE user_id = $2 AND book_id = $3", true, u.ID, b.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
