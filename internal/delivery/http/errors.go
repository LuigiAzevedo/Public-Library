package handler

// HTTP error response message
const (
	timeout            = "request timed out"
	invalidRequestBody = "the request body is invalid or malformed"
)

// Book error response message
const (
	getBook        = "failed to retrieve the book"
	bookNotFound   = "the requested book was not found"
	createBook     = "failed to create the book"
	updateBook     = "failed to update the book"
	deleteBook     = "failed to delete the book"
	searchBook     = "failed to search for books"
	wrongBodyTitle = "invalid title format, it should be a string"
	invalidBookID  = "invalid book ID provided, it should be a positive integer"
)

// User error response message
const (
	getUser       = "failed to retrieve the user"
	createUser    = "failed to create the user"
	updateUser    = "failed to update the user"
	deleteUser    = "failed to delete the user"
	userNotFound  = "the requested user was not found"
	invalidUserID = "invalid user ID provided, it should be a positive integer"
	alreadyExists = "the username or email already exists"
)

// Loan error response message
const (
	borrowBook          = "failed to borrow the book"
	returnBook          = "failed to return the book"
	loanAlreadyReturned = "the loan does not exist or has already been returned"
	loanNotFound        = "no loans found for the user"
	searchUserLoans     = "failed to search for user loans"
	returnBookFirst     = "return the book before borrowing it again"
	bookUnavailable     = "the book is currently unavailable"
)
