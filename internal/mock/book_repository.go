package mock

import (
	"context"
	"strings"
	"time"

	err "github.com/LuigiAzevedo/public-library-v2/internal/database/repository"
	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	ports "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
)

type mockBookRepository struct {
	books []*entity.Book
}

func NewMockBookRepository() ports.BookRepository {
	return &mockBookRepository{
		books: []*entity.Book{
			{
				ID:        1,
				Title:     "Book One",
				Author:    "Author One",
				Amount:    2,
				CreatedAt: time.Now(),
			},
			{
				ID:        2,
				Title:     "Book Two",
				Author:    "Author Two",
				Amount:    0,
				CreatedAt: time.Now(),
			},
		},
	}
}

func (r *mockBookRepository) Get(ctx context.Context, id int) (*entity.Book, error) {
	for _, b := range r.books {
		if b.ID == id {
			return b, nil
		}
	}

	return nil, err.ErrBookNotFound
}

func (r *mockBookRepository) List(ctx context.Context) ([]*entity.Book, error) {
	return r.books, nil
}

func (r *mockBookRepository) Search(ctx context.Context, query string) ([]*entity.Book, error) {
	var result []*entity.Book

	for _, b := range r.books {
		if strings.Contains(strings.ToLower(b.Title), strings.ToLower(query)) {
			result = append(result, b)
		}
	}

	if len(result) == 0 {
		return nil, err.ErrBookNotFound
	}

	return result, nil
}

func (r *mockBookRepository) Create(ctx context.Context, b *entity.Book) (int, error) {
	b.ID = r.books[len(r.books)-1].ID + 1
	b.CreatedAt = time.Now()

	r.books = append(r.books, b)

	return b.ID, nil
}

func (r *mockBookRepository) Update(ctx context.Context, b *entity.Book) error {
	for i, book := range r.books {
		if book.ID == b.ID {
			r.books[i] = b
			return nil
		}
	}

	return err.ErrBookNotFound
}

func (r *mockBookRepository) Delete(ctx context.Context, id int) error {
	for i, book := range r.books {
		if book.ID == id {
			r.books = append(r.books[:i], r.books[i+1:]...)
			return nil
		}
	}

	return err.ErrBookNotFound
}
