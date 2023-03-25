package storage

import (
	"database/sql"
	"time"
)

type OrderItem struct {
	Number     string
	UserID     int
	Status     string
	Accrual    *float64
	UploadedAt time.Time
}

type orderItemRow struct {
	number     string
	userID     int
	status     string
	accrual    sql.NullFloat64
	uploadedAt time.Time
}

type Balance struct {
	Current   float64
	Withdrawn float64
}

type WithdrawalItem struct {
	Order       string
	Sum         float64
	ProcessedAt time.Time
}
