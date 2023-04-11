package ports

import "github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"

type BookUsecase interface {
	GetBook(id int) (*entity.Book, error)
	SearchBook(query string) ([]*entity.Book, error)
	CreateBook(b *entity.Book) (int, error)
	UpdateBook(b *entity.Book) error
	DeleteBook(id int) error
}
