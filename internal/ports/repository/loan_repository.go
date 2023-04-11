package ports

import "github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"

type LoanRepository interface {
	CheckNotReturned(userID, bookID int) error
	Search(userID int) ([]*entity.Loan, error)
	Create(l *entity.Loan) (int, error)
	Update(l *entity.Loan) error
}
