package gophermart

type userAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type orderListItem struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

type orderList []orderListItem

type userBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type userBalanceWithdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type userBalanceWithdrawal struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type userBalanceWithdrawals []userBalanceWithdrawal
