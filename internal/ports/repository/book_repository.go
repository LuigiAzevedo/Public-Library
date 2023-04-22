package ports

import (
	"context"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
)

type BookRepository interface {
	Get(ctx context.Context, id int) (*entity.Book, error)
	List(ctx context.Context) ([]*entity.Book, error)
	Search(ctx context.Context, query string) ([]*entity.Book, error)
	Create(ctx context.Context, b *entity.Book) (int, error)
	Update(ctx context.Context, b *entity.Book) error
	Delete(ctx context.Context, id int) error
}
