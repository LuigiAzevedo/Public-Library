package usecase

import "errors"

// Use Case Errors
var (
	ErrBookNotFound        = errors.New("book not found")
	ErrBookUnavailable     = errors.New("book unavailable at the moment")
	ErrLoanAlreadyReturned = errors.New("loan does't exists or already returned")
	ErrReturnBookFirst     = errors.New("return the book first before borrowing it again")
)
