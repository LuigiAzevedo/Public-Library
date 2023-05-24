package mock

import (
	"context"
	"time"

	err "github.com/LuigiAzevedo/public-library-v2/internal/database/repository"
	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	ports "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
)

type mockUserRepository struct {
	users []*entity.User
}

func NewMockUserRepository() ports.UserRepository {
	return &mockUserRepository{
		users: []*entity.User{
			{
				ID:        1,
				Username:  "UserOne",
				Password:  "PasswordOne",
				Email:     "one@email.com",
				CreatedAt: time.Now(),
			},
			{
				ID:        2,
				Username:  "UserTwo",
				Password:  "PasswordTwo",
				Email:     "two@email.com",
				CreatedAt: time.Now(),
			},
		},
	}
}

func (r *mockUserRepository) Get(ctx context.Context, id int) (*entity.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}

	return nil, err.ErrUserNotFound
}

func (r *mockUserRepository) Create(ctx context.Context, u *entity.User) (int, error) {
	u.ID = r.users[len(r.users)-1].ID + 1
	u.CreatedAt = time.Now()

	r.users = append(r.users, u)

	return u.ID, nil
}

func (r *mockUserRepository) Update(ctx context.Context, u *entity.User) error {
	for i, user := range r.users {
		if user.ID == u.ID {
			r.users[i] = u
			return nil
		}
	}

	return err.ErrUserNotFound
}

func (r *mockUserRepository) Delete(ctx context.Context, id int) error {
	for i, user := range r.users {
		if user.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}

	return err.ErrUserNotFound
}
