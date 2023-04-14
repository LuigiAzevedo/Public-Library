package entity

import (
	"errors"
	"time"
)

type Loan struct {
	ID          int `json:"id"`
	UserID      int `json:"user_id"`
	BookID      int `json:"book_id"`
	Is_returned bool
	CreatedAt   time.Time
}

// NewLoan creates a new loan entity
func NewLoan(userID, bookID int) (*Loan, error) {
	loan := &Loan{
		UserID: userID,
		BookID: bookID,
	}

	err := loan.Validate()
	if err != nil {
		return nil, err
	}

	return loan, nil
}

// Validate validates the loan entity.
func (loan *Loan) Validate() error {
	if loan.UserID <= 0 || loan.BookID <= 0 {
		return errors.New("user ID and book ID can't be empty")
	}

	return nil
}
