package balance

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BalanceHandler struct{}

func NewBalanceHandler() *BalanceHandler {
	return &BalanceHandler{}
}

func (bh *BalanceHandler) GetBalance(c *gin.Context) {
	c.String(http.StatusOK, "GetBalance\n")
}

func (bh *BalanceHandler) BalanceWithdraw(c *gin.Context) {
	c.String(http.StatusOK, "BalanceWithdraw\n")
}

func (bh *BalanceHandler) ListWithdrawals(c *gin.Context) {
	c.String(http.StatusOK, "ListWithdrawals\n")
}
