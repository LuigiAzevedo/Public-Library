package errs

// book errors
var (
	ErrUpdateBook      = "error has occurred while updating the book"
	ErrDeleteBook      = "error has occurred while deleting the book"
	ErrCreateBook      = "error has occurred while creating the book"
	ErrSearchBook      = "error has occurred while searching for books"
	ErrGetBook         = "error has occurred while getting the book"
	ErrBookNotFound    = "book not found"
	ErrInvalidBookID   = "book id should be a positive integer"
	ErrBookUnavailable = "book unavailable at the moment"
)

// user errors
var (
	ErrGetUser       = "error has occurred while getting the user"
	ErrCreateUser    = "error has occurred while creating the user"
	ErrUpdateUser    = "error has occurred while updating the user"
	ErrDeleteUser    = "error has occurred while deleting the user"
	ErrUserNotFound  = "user not found"
	ErrInvalidUserID = "user id should be a positive integer"
	ErrAlreadyExists = "username or email already exists"
)

// loan errors
var (
	ErrBorrowBook          = "an error occurred while borrowing the book"
	ErrReturnBook          = "an error occurred while returning the book"
	ErrLoanAlreadyReturned = "loan does't exists or already returned"
	ErrNoLoansFound        = "user does not have any loan"
	ErrSearchUserLoans     = "an error occurred while searching user loans"
	ErrReturnBookFirst     = "return the book first before borrowing it again"
)

// HTTP error
var (
	ErrTimeout            = "request timed out"
	ErrWrongBodyTitle     = "title should be of the type string"
	ErrInvalidRequestBody = "invalid request body"
)
