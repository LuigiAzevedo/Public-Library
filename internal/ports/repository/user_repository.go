package ports

import "github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"

type UserRepository interface {
	Get(id int) (*entity.User, error)
	Create(u *entity.User) (int, error)
	Update(u *entity.User) error
	Delete(id int) error
}
