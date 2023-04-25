package entity

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	tests := map[string]struct {
		username string
		password string
		email    string
		want     error
	}{
		"OK": {
			username: "luigi",
			password: "secret",
			email:    "luigi@email.com",
			want:     nil,
		},
		"Empty Fields": {
			username: "",
			password: "",
			email:    "",
			want:     errors.New("username, password and email can't be empty"),
		},
		"Fields With Spaces": {
			username: "User Name",
			password: "Pass Word",
			email:    "luigi@email.com",
			want:     errors.New("username and password can't have spaces"),
		},
		"Short Password": {
			username: "luigi",
			password: "short",
			email:    "luigi@email.com",
			want:     errors.New("password shorter than 6 characters"),
		},
		"Long Password": {
			username: "luigi",
			password: "ReallyLongPasswordReallyLongPasswordReallyLongPasswordReallyLongPassword2",
			email:    "luigi@email.com",
			want:     errors.New("password longer than 72 characters"),
		},
		"Invalid Email": {
			username: "luigi",
			password: "secret",
			email:    "luigiEmail.com",
			want:     errors.New("invalid email address"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			u, err := NewUser(tc.username, tc.password, tc.email)

			assert.Equal(t, tc.want, err)

			if name == "OK" {
				assert.Equal(t, tc.username, u.Username)
				assert.Equal(t, tc.password, u.Password)
				assert.Equal(t, tc.email, u.Email)
			}
		})
	}
}
