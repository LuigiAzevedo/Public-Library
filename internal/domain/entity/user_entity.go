package entity

import (
	"net/mail"
	"strings"
	"time"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user entity
func NewUser(username, password, email string) (*User, error) {
	user := &User{
		Username:  username,
		Password:  password,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate validates the user entity.
func (user *User) Validate() error {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return ErrEmptyUserField
	}

	if strings.ContainsAny(user.Username, " \t\r\n") || strings.ContainsAny(user.Password, " \t\r\n") {
		return ErrFieldWithSpaces
	}

	if len(user.Password) < 6 {
		return ErrShortPassword
	}

	if len(user.Password) > 72 {
		return ErrLongPassword
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return ErrInvalidEmail
	}

	return nil
}
