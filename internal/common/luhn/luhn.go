package luhn

import (
	"github.com/ShiraazMoollatjie/goluhn"
)

func CheckNum(orderNum string) bool {
	return goluhn.Validate(orderNum) == nil
}
