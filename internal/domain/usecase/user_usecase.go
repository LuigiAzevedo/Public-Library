package usecase

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	r "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
	u "github.com/LuigiAzevedo/public-library-v2/internal/ports/usecase"
)

type userUseCase struct {
	userRepo r.UserRepository
}

// NewUserUseCase creates a new instance of userUseCase
func NewUserUseCase(repository r.UserRepository) u.UserUsecase {
	return &userUseCase{
		userRepo: repository,
	}
}

func (s *userUseCase) GetUser(ctx context.Context, id int) (*entity.User, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userUseCase) CreateUser(ctx context.Context, u *entity.User) (int, error) {
	user, err := entity.NewUser(u.Username, u.Password, u.Email)
	if err != nil {
		return 0, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user.Password = string(hashedPassword)

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *userUseCase) UpdateUser(ctx context.Context, u *entity.User) error {
	u.UpdatedAt = time.Now()

	err := u.Validate()
	if err != nil {
		return err
	}

	err = s.userRepo.Update(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

func (s *userUseCase) DeleteUser(ctx context.Context, id int) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
