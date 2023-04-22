package ports

import (
	"context"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
)

type UserRepository interface {
	Get(ctx context.Context, id int) (*entity.User, error)
	Create(ctx context.Context, u *entity.User) (int, error)
	Update(ctx context.Context, u *entity.User) error
	Delete(ctx context.Context, id int) error
}
