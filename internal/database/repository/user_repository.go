package repository

import (
	"database/sql"
	"errors"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	r "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
)

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
func (r *userRepository) Get(id int) (*entity.User, error) {
	stmt, err := r.db.Prepare("SELECT * FROM users WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	u := &entity.User{}

	row := stmt.QueryRow(id)

	var updatedAt sql.NullTime
	err = row.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &updatedAt, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	// check if updated_at is NULL before scanning it
	if updatedAt.Valid {
		u.UpdatedAt = updatedAt.Time
	}

	return u, nil
}

// Create creates a new user
func (r *userRepository) Create(u *entity.User) (int, error) {
	stmt, err := r.db.Prepare("INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(u.Username, u.Password, u.Email).Scan(&u.ID)
	if err != nil {
		return 0, err
	}

	return u.ID, nil
}

// Update updates an user
func (r *userRepository) Update(u *entity.User) error {
	stmt, err := r.db.Prepare("UPDATE users SET username = $1, password = $2, email = $3, updated_at = $4 WHERE id = $5")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Username, u.Password, u.Email, u.UpdatedAt, u.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("0 rows affected")
	}

	return nil
}

// Delete deletes an user by id
func (r *userRepository) Delete(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("0 rows affected")
	}

	return nil
}
