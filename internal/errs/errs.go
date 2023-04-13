package errs

// book errors
var (
	ErrUpdateBook = "error has occurred while updating the book"
	ErrDeleteBook = "error has occurred while deleting the book"
	ErrCreateBook = "error has occurred while creating the book"
	ErrSearchBook = "error has occurred while searching for books"
	ErrGetBook    = "error has occurred while getting the book"
)

// user errors
var (
	ErrGetUser    = "error has occurred while getting the user"
	ErrCreateUser = "error has occurred while creating the user"
	ErrUpdateUser = "error has occurred while updating the user"
	ErrDeleteUser = "error has occurred while deleting the user"
)

// loan errors
var (
	ErrBorrowBook = "an error occurred while borrowing the book"
)
