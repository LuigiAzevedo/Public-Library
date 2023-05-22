package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	"github.com/LuigiAzevedo/public-library-v2/internal/errs"
	r "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
	u "github.com/LuigiAzevedo/public-library-v2/internal/ports/usecase"
)

type loanUseCase struct {
	loanRepo r.LoanRepository
	userRepo r.UserRepository
	bookRepo r.BookRepository
}

// NewLoanUseCase creates a new instance of loanUseCase
func NewLoanUseCase(loan r.LoanRepository, user r.UserRepository, book r.BookRepository) u.LoanUsecase {
	return &loanUseCase{
		loanRepo: loan,
		userRepo: user,
		bookRepo: book,
	}
}

func (s *loanUseCase) BorrowBook(ctx context.Context, userID, bookID int) error {
	exists, err := s.loanRepo.CheckNotReturned(ctx, userID, bookID)
	if err != nil {
		return fmt.Errorf("%s: %w", errs.ErrBorrowBook, err)
	}
	if exists {
		return errors.New(errs.ErrReturnBookFirst)
	}

	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", errs.ErrGetUser, err)
	}

	book, err := s.bookRepo.Get(ctx, bookID)
	if err != nil {
		return fmt.Errorf("%s: %w", errs.ErrGetBook, err)
	}

	book.Amount -= 1
	if book.Amount < 0 {
		return errors.New(errs.ErrBookUnavailable)
	}

	err = s.loanRepo.BorrowTransaction(ctx, user, book)
	if err != nil {
		return fmt.Errorf("%s: %w", errs.ErrBorrowBook, err)
	}

	return nil
}

func (s *loanUseCase) ReturnBook(ctx context.Context, userID, bookID int) error {
	exists, err := s.loanRepo.CheckNotReturned(ctx, userID, bookID)
	if err != nil {
		return fmt.Errorf("%s: %w", errs.ErrReturnBook, err)
	}
	if !exists {
		return errors.New(errs.ErrLoanAlreadyReturned)
	}

	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", errs.ErrGetUser, err)
	}

	book, err := s.bookRepo.Get(ctx, bookID)
	if err != nil {
		return fmt.Errorf("%s: %w", errs.ErrGetBook, err)
	}

	book.Amount += 1

	err = s.loanRepo.ReturnTransaction(ctx, user, book)
	if err != nil {
		return fmt.Errorf("%s: %w", errs.ErrReturnBook, err)
	}

	return nil
}

func (s *loanUseCase) SearchUserLoans(ctx context.Context, userID int) ([]*entity.Loan, error) {
	loans, err := s.loanRepo.Search(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errs.ErrSearchUserLoans, err)
	}

	return loans, nil
}
