package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	"github.com/LuigiAzevedo/public-library-v2/internal/errs"
	r "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
	u "github.com/LuigiAzevedo/public-library-v2/internal/ports/usecase"
)

type bookUseCase struct {
	bookRepo r.BookRepository
}

// NewBookUseCase creates a new instance of bookUseCase
func NewBookUseCase(repository r.BookRepository) u.BookUsecase {
	return &bookUseCase{
		bookRepo: repository,
	}
}

func (s *bookUseCase) GetBook(ctx context.Context, id int) (*entity.Book, error) {
	book, err := s.bookRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errs.ErrGetBook, err)
	}

	return book, nil
}

func (s *bookUseCase) SearchBooks(ctx context.Context, query string) ([]*entity.Book, error) {
	books, err := s.bookRepo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errs.ErrSearchBook, err)
	}

	return books, nil
}

func (s *bookUseCase) ListBooks(ctx context.Context) ([]*entity.Book, error) {
	books, err := s.bookRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errs.ErrSearchBook, err)
	}

	return books, nil
}

func (s *bookUseCase) CreateBook(ctx context.Context, b *entity.Book) (int, error) {
	book, err := entity.NewBook(b.Title, b.Author, b.Amount)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errs.ErrCreateBook, err)
	}

	id, err := s.bookRepo.Create(ctx, book)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errs.ErrCreateBook, err)
	}

	return id, nil
}

func (s *bookUseCase) UpdateBook(ctx context.Context, b *entity.Book) error {
	b.UpdatedAt = time.Now()

	err := b.Validate()
	if err != nil {
		return fmt.Errorf("%s: %w", errs.ErrUpdateBook, err)
	}

	err = s.bookRepo.Update(ctx, b)
	if err != nil {
		return fmt.Errorf("%s: %w", errs.ErrUpdateBook, err)
	}

	return nil
}

func (s *bookUseCase) DeleteBook(ctx context.Context, id int) error {
	err := s.bookRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", errs.ErrDeleteBook, err)
	}

	return nil
}
