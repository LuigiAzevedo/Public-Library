package entity

import (
	"time"
)

type Book struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Amount    int    `json:"amount"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewBook creates a new book entity
func NewBook(title, author string, amount int) (*Book, error) {
	book := &Book{
		Title:     title,
		Author:    author,
		Amount:    amount,
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}

	if err := book.Validate(); err != nil {
		return nil, err
	}

	return book, nil
}

// Validate validates the book entity.
func (book *Book) Validate() error {
	if book.Title == "" || book.Author == "" || book.Amount <= 0 {
		return ErrInvalidBook
	}

	return nil
}
