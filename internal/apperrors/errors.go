package apperrors

import "errors"

var (
	ErrAccountExists            = errors.New("account already exists")
	ErrAccountNotFound          = errors.New("account not found")
	ErrInvalidAccountID         = errors.New("invalid account ID")
	ErrInternalServerError      = errors.New("internal server error")
	ErrTransactionAlreadyExists = errors.New("transaction already exists")
	ErrInvalidOperationType     = errors.New("invalid operation type")
	ErrAccountDoesNotExist      = errors.New("account does not exist")
)
