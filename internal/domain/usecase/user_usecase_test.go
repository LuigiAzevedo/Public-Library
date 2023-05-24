package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	"github.com/LuigiAzevedo/public-library-v2/internal/mock"
)

func TestGetUser(t *testing.T) {
	repo := mock.NewMockUserRepository()
	uc := NewUserUseCase(repo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		u, err := uc.GetUser(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, "UserOne", u.Username)
	})
	t.Run("Not Found", func(t *testing.T) {
		_, err := uc.GetUser(ctx, 0)
		assert.Error(t, err)
	})
}

func TestCreateUser(t *testing.T) {
	repo := mock.NewMockUserRepository()
	uc := NewUserUseCase(repo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		u := &entity.User{
			Username: "UserThree",
			Password: "PasswordThree",
			Email:    "three@email.com",
		}

		id, err := uc.CreateUser(ctx, u)
		assert.NoError(t, err)
		assert.Equal(t, 3, id)
	})
	t.Run("Invalid User", func(t *testing.T) {
		u := &entity.User{}

		_, err := uc.CreateUser(ctx, u)
		assert.Error(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	repo := mock.NewMockUserRepository()
	uc := NewUserUseCase(repo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		u := &entity.User{
			ID:       2,
			Username: "UserFive",
			Password: "PasswordFive",
			Email:    "five@email.com",
		}

		err := uc.UpdateUser(ctx, u)
		assert.NoError(t, err)
	})
	t.Run("Not Found", func(t *testing.T) {
		u := &entity.User{
			ID:       5,
			Username: "UserFive",
			Password: "PasswordFive",
			Email:    "five@email.com",
		}

		err := uc.UpdateUser(ctx, u)
		assert.Error(t, err)
	})
}

func TestDeleteUser(t *testing.T) {
	repo := mock.NewMockUserRepository()
	uc := NewUserUseCase(repo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		err := uc.DeleteUser(ctx, 2)
		assert.NoError(t, err)
	})
	t.Run("Not Found", func(t *testing.T) {
		err := uc.DeleteUser(ctx, 5)
		assert.Error(t, err)
	})
}
