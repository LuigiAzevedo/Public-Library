package usecase

import (
	"errors"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	"github.com/LuigiAzevedo/public-library-v2/internal/errs"
	r "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
	u "github.com/LuigiAzevedo/public-library-v2/internal/ports/usecase"
)

type loanService struct {
	loanRepo r.LoanRepository
	userRepo r.UserRepository
	bookRepo r.BookRepository
}

// NewLoanService creates a new instance of userService
func NewLoanService(repository r.LoanRepository) u.LoanUsecase {
	return &loanService{
		loanRepo: repository,
	}
}

func (s *loanService) BorrowBook(userID, bookID int) error {
	exists, err := s.loanRepo.CheckNotReturned(userID, bookID)
	if err != nil {
		return errs.ErrorWrapper(errs.ErrBorrowBook, err)
	}
	if exists {
		return errors.New("return the book first before borrowing it again")
	}

	user, err := s.userRepo.Get(userID)
	if err != nil {
		return errs.ErrorWrapper(errs.ErrGetUser, err)
	}

	book, err := s.bookRepo.Get(bookID)
	if err != nil {
		return errs.ErrorWrapper(errs.ErrGetBook, err)
	}

	book.Amount -= 1
	if book.Amount < 0 {
		return errors.New("book unavailable at the moment")
	}

	err = s.loanRepo.BorrowTransaction(user, book)
	if err != nil {
		return errs.ErrorWrapper(errs.ErrBorrowBook, err)
	}

	return nil
}

func (s *loanService) ReturnBook(userID, bookID int) error {
	exists, err := s.loanRepo.CheckNotReturned(userID, bookID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("loan does't exists or already returned")
	}

	user, err := s.userRepo.Get(userID)
	if err != nil {
		return errs.ErrorWrapper(errs.ErrGetUser, err)
	}

	book, err := s.bookRepo.Get(bookID)
	if err != nil {
		return errs.ErrorWrapper(errs.ErrGetBook, err)
	}

	book.Amount += 1

	err = s.loanRepo.ReturnTransaction(user, book)
	if err != nil {
		return errs.ErrorWrapper("an error occurred while returning the book", err)
	}

	return nil
}

func (s *loanService) SearchUserLoans(userID int) ([]*entity.Loan, error) {
	loans, err := s.loanRepo.Search(userID)
	if err != nil {
		return nil, errs.ErrorWrapper("an error occurred while searching loans", err)
	}

	return loans, nil
}
