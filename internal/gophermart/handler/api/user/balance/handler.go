package balance

import (
	"errors"
	"net/http"
	"time"

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

type WithdrawalItemResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
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
	dbItems, err := bh.stg.ListWithdrawals(auth.GetUserId(c))
	if err != nil {
		if errors.Is(storage.ErrNoWithdrawals, err) {
			_http.CreateStatusResponse(c, http.StatusNoContent)
			return
		}

		_http.CreateStatusResponse(c, http.StatusInternalServerError)
		return
	}

	respItems := make([]WithdrawalItemResponse, 0, len(dbItems))
	for _, item := range dbItems {
		respItems = append(respItems, WithdrawalItemResponse{
			Order:       item.Order,
			Sum:         item.Sum,
			ProcessedAt: item.ProcessedAt,
		})
	}

	c.JSON(http.StatusOK, respItems)
}
