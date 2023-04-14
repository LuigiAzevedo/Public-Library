package entity

import (
	"errors"
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

	err := user.Validate()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Validate validates the user entity.
func (user *User) Validate() error {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return errors.New("username, password and email can't be empty")
	}

	if strings.ContainsAny(user.Username, " \t\r\n") || strings.ContainsAny(user.Password, " \t\r\n") {
		return errors.New("username and password can't have spaces")
	}

	if len(user.Password) < 6 {
		return errors.New("password shorter than 6 characters")
	}

	if len(user.Password) > 72 {
		return errors.New("password longer than 72 characters")
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return errors.New("invalid email address")
	}

	return nil
}
