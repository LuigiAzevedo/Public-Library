package ports

import (
	"context"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
)

type BookUsecase interface {
	GetBook(ctx context.Context, id int) (*entity.Book, error)
	ListBooks(ctx context.Context) ([]*entity.Book, error)
	SearchBooks(ctx context.Context, query string) ([]*entity.Book, error)
	CreateBook(ctx context.Context, b *entity.Book) (int, error)
	UpdateBook(ctx context.Context, b *entity.Book) error
	DeleteBook(ctx context.Context, id int) error
}
