package gophermart

import "errors"

var (
	// Common
	ErrUnauthorized  = errors.New("unauthorized request")
	ErrBadRequest    = errors.New("bad request")
	ErrInternalError = errors.New("internal error")
	// User
	ErrUserAlreadyExists = errors.New("user already exists")
	// Balance
	ErrNoBalanceWithdrawals           = errors.New("no balance withdrawals")
	ErrBalanceWithdrawWrongFormat     = errors.New("balance withdraw wrong format")
	ErrBalanceWithdrawPaymentRequired = errors.New("balance withdraw payment required")
	// Orders
	ErrOrderWrongFormat     = errors.New("order wrong format")
	ErrOrderAlreadyExists   = errors.New("order already exists")
	ErrOrderAlreadyAccepted = errors.New("order already accepted")
	ErrNoOrders             = errors.New("no orders")
)
