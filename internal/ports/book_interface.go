package ports

import "github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"

type BookUsecase interface {
	Get(id int) (*entity.Book, error)
	SearchAndList(query string) ([]*entity.Book, error)
	Create(b *entity.Book) (string, error)
	Update(b *entity.Book) error
	Delete(id int) error
}
