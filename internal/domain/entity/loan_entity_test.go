package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoan(t *testing.T) {
	tests := map[string]struct {
		userID int
		bookID int
		want   error
	}{
		"OK": {
			userID: 1,
			bookID: 1,
			want:   nil,
		},
		"Invalid UserID": {
			userID: 0,
			bookID: 1,
			want:   ErrInvalidLoan,
		},
		"Invalid BookID": {
			userID: 1,
			bookID: 0,
			want:   ErrInvalidLoan,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			l, err := NewLoan(tc.userID, tc.bookID)

			assert.Equal(t, tc.want, err)

			if name == "OK" {
				assert.Equal(t, tc.userID, l.UserID)
				assert.Equal(t, tc.bookID, l.BookID)
			}
		})
	}
}
