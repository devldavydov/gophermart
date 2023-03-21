package luhn

import (
	"strconv"

	"github.com/theplant/luhn"
)

func CheckNum(orderNum string) bool {
	order, err := strconv.Atoi(orderNum)
	if err != nil {
		return false
	}
	return luhn.Valid(order)
}
