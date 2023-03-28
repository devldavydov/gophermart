package gophermart

import "errors"

var (
	// Common
	ErrUnauthorized  = errors.New("unauthorized request")
	ErrBadRequest    = errors.New("bad request")
	ErrInternalError = errors.New("internal error")
	// User
	ErrUserAlreadyExists = errors.New("user already exists")
)
