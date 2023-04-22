package ports

import (
	"context"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
)

type LoanRepository interface {
	CheckNotReturned(ctx context.Context, userID, bookID int) (bool, error)
	Search(ctx context.Context, userID int) ([]*entity.Loan, error)
	BorrowTransaction(ctx context.Context, u *entity.User, b *entity.Book) error
	ReturnTransaction(ctx context.Context, u *entity.User, b *entity.Book) error
}
