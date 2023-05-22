package usecase

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	"github.com/LuigiAzevedo/public-library-v2/internal/errs"
	"github.com/LuigiAzevedo/public-library-v2/internal/mock"
)

func TestGetBook(t *testing.T) {
	repo := mock.NewMockBookRepository()
	uc := NewBookUseCase(repo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		b, err := uc.GetBook(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, "Book One", b.Title)
	})
	t.Run("Not Found", func(t *testing.T) {
		_, err := uc.GetBook(ctx, 0)
		assert.Equal(t, errs.ErrBookNotFound, errors.Unwrap(err).Error())
	})
}

func TestSearchBooks(t *testing.T) {
	repo := mock.NewMockBookRepository()
	uc := NewBookUseCase(repo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		b, err := uc.SearchBooks(ctx, "two")
		assert.NoError(t, err)
		assert.Equal(t, "Book Two", b[0].Title)
	})
	t.Run("Not Found", func(t *testing.T) {
		_, err := uc.SearchBooks(ctx, "five")
		assert.Equal(t, errs.ErrBookNotFound, errors.Unwrap(err).Error())
	})
}

func TestListBooks(t *testing.T) {
	repo := mock.NewMockBookRepository()
	uc := NewBookUseCase(repo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		b, err := uc.ListBooks(ctx)
		assert.NoError(t, err)

		for _, book := range b {
			assert.NotEmpty(t, book)
		}
	})
}

func TestCreateBook(t *testing.T) {
	repo := mock.NewMockBookRepository()
	uc := NewBookUseCase(repo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		b := &entity.Book{
			Title:  "Book Three",
			Author: "Author Three",
			Amount: 5,
		}

		id, err := uc.CreateBook(ctx, b)
		assert.NoError(t, err)
		assert.Equal(t, 3, id)
	})
}

func TestUpdateBook(t *testing.T) {
	repo := mock.NewMockBookRepository()
	uc := NewBookUseCase(repo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		b := &entity.Book{
			ID:     2,
			Title:  "Book Three",
			Author: "Author Three",
			Amount: 5,
		}

		err := uc.UpdateBook(ctx, b)
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		b := &entity.Book{
			ID:     5,
			Title:  "Book Three",
			Author: "Author Three",
			Amount: 5,
		}

		err := uc.UpdateBook(ctx, b)
		assert.Equal(t, errs.ErrBookNotFound, errors.Unwrap(err).Error())
	})
}

func TestDeleteBook(t *testing.T) {
	repo := mock.NewMockBookRepository()
	uc := NewBookUseCase(repo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		err := uc.DeleteBook(ctx, 2)
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		err := uc.DeleteBook(ctx, 5)
		assert.Equal(t, errs.ErrBookNotFound, errors.Unwrap(err).Error())
	})
}
