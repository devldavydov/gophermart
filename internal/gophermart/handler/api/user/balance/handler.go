package balance

import (
	"net/http"

	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BalanceHandler struct {
	stg    storage.Storage
	logger *logrus.Logger
}

func NewBalanceHandler(stg storage.Storage, logger *logrus.Logger) *BalanceHandler {
	return &BalanceHandler{stg: stg, logger: logger}
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
