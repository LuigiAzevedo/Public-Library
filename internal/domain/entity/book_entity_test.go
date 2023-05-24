package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBook(t *testing.T) {
	tests := map[string]struct {
		title  string
		author string
		amount int
		want   error
	}{
		"OK": {
			title:  "Let's Go Further!",
			author: "Alex Edwards",
			amount: 5,
			want:   nil,
		},
		"Empty Fields": {
			title:  "",
			author: "",
			amount: 5,
			want:   ErrInvalidBook,
		},
		"Invalid Amount": {
			title:  "Let's Go Further!",
			author: "Alex Edwards",
			amount: 0,
			want:   ErrInvalidBook,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := NewBook(tc.title, tc.author, tc.amount)

			assert.Equal(t, tc.want, err)

			if name == "OK" {
				assert.Equal(t, tc.title, b.Title)
				assert.Equal(t, tc.author, b.Author)
				assert.Equal(t, tc.amount, b.Amount)
			}
		})
	}
}
