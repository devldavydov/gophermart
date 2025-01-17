package gophermart

type UserBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type UserBalanceWithdrawal struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type UserBalanceWithdrawals []UserBalanceWithdrawal

type OrderItem struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

type OrderList []OrderItem

type userAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type userBalanceWithdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
