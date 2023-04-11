package ports

import "github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"

type BookRepository interface {
	Get(id int) (*entity.Book, error)
	List() ([]*entity.Book, error)
	Search(query string) ([]*entity.Book, error)
	Create(b *entity.Book) (int, error)
	Update(b *entity.Book) error
	Delete(id int) error
}
