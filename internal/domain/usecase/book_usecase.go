package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"

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
		return nil, errors.Wrap(err, errs.ErrGetBook)
	}

	return book, nil
}

func (s *bookUseCase) SearchBooks(ctx context.Context, query string) ([]*entity.Book, error) {
	books, err := s.bookRepo.Search(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, errs.ErrSearchBook)
	}

	return books, nil
}

func (s *bookUseCase) ListBooks(ctx context.Context) ([]*entity.Book, error) {
	books, err := s.bookRepo.List(ctx)
	if err != nil {
		return nil, errors.Wrap(err, errs.ErrSearchBook)
	}

	return books, nil
}

func (s *bookUseCase) CreateBook(ctx context.Context, b *entity.Book) (int, error) {
	book, err := entity.NewBook(b.Title, b.Author, b.Amount)
	if err != nil {
		return 0, errors.Wrap(err, errs.ErrCreateBook)
	}

	id, err := s.bookRepo.Create(ctx, book)
	if err != nil {
		return 0, errors.Wrap(err, errs.ErrCreateBook)
	}

	return id, nil
}

func (s *bookUseCase) UpdateBook(ctx context.Context, b *entity.Book) error {
	b.UpdatedAt = time.Now()

	err := b.Validate()
	if err != nil {
		return errors.Wrap(err, errs.ErrUpdateBook)
	}

	err = s.bookRepo.Update(ctx, b)
	if err != nil {
		return errors.Wrap(err, errs.ErrUpdateBook)
	}

	return nil
}

func (s *bookUseCase) DeleteBook(ctx context.Context, id int) error {
	err := s.bookRepo.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, errs.ErrDeleteBook)
	}

	return nil
}
