package storage

import (
	"database/sql"
	"time"
)

type OrderItem struct {
	Number     string
	Status     string
	Accrual    *int32
	UploadedAt time.Time
}

type orderItemRow struct {
	number     string
	status     string
	accrual    sql.NullInt32
	uploadedAt time.Time
}

type Balance struct {
	Current   float64
	Withdrawn float64
}
