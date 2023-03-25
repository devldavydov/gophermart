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
