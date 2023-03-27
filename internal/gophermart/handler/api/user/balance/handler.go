package balance

import (
	"errors"
	"net/http"
	"time"

	_http "github.com/devldavydov/gophermart/internal/common/http"
	"github.com/devldavydov/gophermart/internal/common/luhn"
	"github.com/devldavydov/gophermart/internal/gophermart/auth"
	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const _notEnougnBalance = "Not enough balance"

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawalItemResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type BalanceWithdrawReq struct {
	Order string  `json:"order" binding:"required"`
	Sum   float64 `json:"sum" binding:"required"`
}

type BalanceHandler struct {
	stg    storage.Storage
	logger *logrus.Logger
}

func NewBalanceHandler(stg storage.Storage, logger *logrus.Logger) *BalanceHandler {
	return &BalanceHandler{stg: stg, logger: logger}
}

func (bh *BalanceHandler) GetBalance(c *gin.Context) {
	dbBalance, err := bh.stg.GetBalance(c.Request.Context(), auth.GetUserID(c))
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
	var req BalanceWithdrawReq

	if err := c.ShouldBindJSON(&req); err != nil {
		_http.CreateStatusResponse(c, http.StatusBadRequest)
		return
	}

	if !luhn.CheckNum(req.Order) {
		_http.CreateStatusResponse(c, http.StatusUnprocessableEntity)
		return
	}

	err := bh.stg.BalanceWithdraw(c.Request.Context(), auth.GetUserID(c), req.Order, req.Sum)
	if err != nil {
		if errors.Is(storage.ErrNotEnoughBalance, err) {
			c.String(http.StatusPaymentRequired, _notEnougnBalance)
			return
		}

		_http.CreateStatusResponse(c, http.StatusInternalServerError)
		return
	}

	_http.CreateStatusResponse(c, http.StatusOK)
}

func (bh *BalanceHandler) ListWithdrawals(c *gin.Context) {
	dbItems, err := bh.stg.ListWithdrawals(c.Request.Context(), auth.GetUserID(c))
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
