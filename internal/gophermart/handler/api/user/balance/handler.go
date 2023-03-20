package balance

import (
	"net/http"

	_http "github.com/devldavydov/gophermart/internal/common/http"
	"github.com/devldavydov/gophermart/internal/gophermart/auth"
	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type BalanceHandler struct {
	stg    storage.Storage
	logger *logrus.Logger
}

func NewBalanceHandler(stg storage.Storage, logger *logrus.Logger) *BalanceHandler {
	return &BalanceHandler{stg: stg, logger: logger}
}

func (bh *BalanceHandler) GetBalance(c *gin.Context) {
	dbBalance, err := bh.stg.GetBalance(auth.GetUserId(c))
	if err != nil {
		_http.CreateStatusResponse(c, http.StatusInternalServerError)
		return
	}

	respBalance := BalanceResponse{
		Current:   dbBalance.Current,
		Withdrawn: dbBalance.Withdrawn,
	}
	c.JSON(http.StatusOK, respBalance)
}

func (bh *BalanceHandler) BalanceWithdraw(c *gin.Context) {
	c.String(http.StatusOK, "BalanceWithdraw\n")
}

func (bh *BalanceHandler) ListWithdrawals(c *gin.Context) {
	c.String(http.StatusOK, "ListWithdrawals\n")
}
