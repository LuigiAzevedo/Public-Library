package ports

import (
	"context"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
)

type LoanUsecase interface {
	BorrowBook(ctx context.Context, userID, bookID int) error
	ReturnBook(ctx context.Context, userID, bookID int) error
	SearchUserLoans(ctx context.Context, userID int) ([]*entity.Loan, error)
}
