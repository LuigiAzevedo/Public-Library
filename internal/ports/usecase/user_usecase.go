package ports

import "github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"

type UserUsecase interface {
	GetUser(id int) (*entity.User, error)
	CreateUser(u *entity.User) (int, error)
	UpdateUser(u *entity.User) error
	DeleteUser(id int) error
}
