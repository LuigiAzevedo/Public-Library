package entity

import (
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
			want:     ErrEmptyUserField,
		},
		"Fields With Spaces": {
			username: "User Name",
			password: "Pass Word",
			email:    "luigi@email.com",
			want:     ErrFieldWithSpaces,
		},
		"Short Password": {
			username: "luigi",
			password: "short",
			email:    "luigi@email.com",
			want:     ErrShortPassword,
		},
		"Long Password": {
			username: "luigi",
			password: "ReallyLongPasswordReallyLongPasswordReallyLongPasswordReallyLongPassword2",
			email:    "luigi@email.com",
			want:     ErrLongPassword,
		},
		"Invalid Email": {
			username: "luigi",
			password: "secret",
			email:    "luigiEmail.com",
			want:     ErrInvalidEmail,
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
