package mock

import (
	"context"
	"time"

	err "github.com/LuigiAzevedo/public-library-v2/internal/database/repository"
	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	ports "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
)

type mockLoanRepository struct {
	loans []*entity.Loan
}

func NewMockLoanRepository() ports.LoanRepository {
	return &mockLoanRepository{
		loans: []*entity.Loan{
			{
				ID:          1,
				UserID:      1,
				BookID:      1,
				Is_returned: true,
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
				UserID:      2,
				BookID:      2,
				Is_returned: false,
				CreatedAt:   time.Now(),
			},
		},
	}
}

func (r *mockLoanRepository) CheckNotReturned(ctx context.Context, userID int, bookID int) (bool, error) {
	for _, l := range r.loans {
		if l.UserID == userID && l.BookID == bookID && !l.Is_returned {
			return true, nil
		}
	}

	return false, nil
}

func (r *mockLoanRepository) Search(ctx context.Context, userID int) ([]*entity.Loan, error) {
	var loans []*entity.Loan
	for _, l := range r.loans {
		if l.UserID == userID {
			loans = append(loans, l)
		}
	}

	if len(loans) == 0 {
		return nil, err.ErrLoanNotFound
	}

	return loans, nil
}

func (r *mockLoanRepository) BorrowTransaction(ctx context.Context, u *entity.User, b *entity.Book) error {
	loan := &entity.Loan{
		ID:        4,
		UserID:    u.ID,
		BookID:    b.ID,
		CreatedAt: time.Now(),
	}

	r.loans = append(r.loans, loan)

	return nil
}

func (r *mockLoanRepository) ReturnTransaction(ctx context.Context, u *entity.User, b *entity.Book) error {
	for _, l := range r.loans {
		if l.UserID == u.ID && l.BookID == b.ID && !l.Is_returned {
			r.loans[l.ID-1].Is_returned = true
			return nil
		}
	}

	return err.ErrLoanNotFound
}
