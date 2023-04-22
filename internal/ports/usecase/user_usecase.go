package ports

import (
	"context"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
)

type UserUsecase interface {
	GetUser(ctx context.Context, id int) (*entity.User, error)
	CreateUser(ctx context.Context, u *entity.User) (int, error)
	UpdateUser(ctx context.Context, u *entity.User) error
	DeleteUser(ctx context.Context, id int) error
}
