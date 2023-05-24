package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	r "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
)

const duplicatedKeyValueCode = "23505"

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) r.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Get gets user info by id
func (r *userRepository) Get(ctx context.Context, id int) (*entity.User, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT * FROM users WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}
	defer stmt.Close()

	u := &entity.User{}

	row := stmt.QueryRowContext(ctx, id)

	var updatedAt sql.NullTime
	err = row.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &updatedAt, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		} else {
			return nil, fmt.Errorf("%s: %w", ErrScanData, err)
		}
	}

	// check if updatedAt is not NULL
	if updatedAt.Valid {
		u.UpdatedAt = updatedAt.Time
	}

	return u, nil
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, u *entity.User) (int, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.Username, u.Password, u.Email).Scan(&u.ID)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == duplicatedKeyValueCode {
			return 0, ErrAlreadyExists
		} else {
			return 0, fmt.Errorf("%s: %w", ErrExecuteQuery, err)
		}
	}

	return u.ID, nil
}

// Update updates an user
func (r *userRepository) Update(ctx context.Context, u *entity.User) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE users SET username = $1, password = $2, email = $3, updated_at = NOW() WHERE id = $4")
	if err != nil {
		return fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, u.Username, u.Password, u.Email, u.ID)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == duplicatedKeyValueCode {
			return ErrAlreadyExists
		} else {
			return fmt.Errorf("%s: %w", ErrExecuteStatement, err)
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrRetrieveRows, err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Delete deletes an user by id
func (r *userRepository) Delete(ctx context.Context, id int) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM users WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrExecuteStatement, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrRetrieveRows, err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
