package gophermart

type UserBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type userAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
