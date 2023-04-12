package ports

import "github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"

type LoanRepository interface {
	CheckNotReturned(userID, bookID int) (bool, error)
	Search(userID int) ([]*entity.Loan, error)
	BorrowTransaction(*entity.User, *entity.Book) error
	ReturnTransaction(*entity.User, *entity.Book) error
}
