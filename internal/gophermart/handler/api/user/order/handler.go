package order

import (
	"errors"
	"io"
	"net/http"
	"time"

	_http "github.com/devldavydov/gophermart/internal/common/http"
	"github.com/devldavydov/gophermart/internal/common/luhn"
	"github.com/devldavydov/gophermart/internal/gophermart/auth"
	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	_orderAlreadyExists         = "Order already exists"
	_orderAlreadyExistsFromUser = "Order already exists from user"
	_orderAccepted              = "Order accepted"
)

type OrderItemsResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type OrderHandler struct {
	stg    storage.Storage
	logger *logrus.Logger
}

func NewOrderHandler(stg storage.Storage, logger *logrus.Logger) *OrderHandler {
	return &OrderHandler{stg: stg, logger: logger}
}

func (oh *OrderHandler) AddOrder(c *gin.Context) {
	if !_http.CheckRequestContentType(c.Request.Header, "text/plain") {
		_http.CreateStatusResponse(c, http.StatusBadRequest)
		return
	}

	orderNumBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		_http.CreateStatusResponse(c, http.StatusBadRequest)
		return
	}

	orderNum := string(orderNumBytes)
	if !luhn.CheckNum(orderNum) {
		_http.CreateStatusResponse(c, http.StatusUnprocessableEntity)
		return
	}

	err = oh.stg.AddOrder(c.Request.Context(), auth.GetUserID(c), orderNum)
	if err != nil {
		if errors.Is(storage.ErrOrderAlreadyExists, err) {
			c.String(http.StatusConflict, _orderAlreadyExists)
			return
		}

		if errors.Is(storage.ErrOrderAlreadyExistsFromUser, err) {
			c.String(http.StatusOK, _orderAlreadyExistsFromUser)
			return
		}

		_http.CreateStatusResponse(c, http.StatusInternalServerError)
		return
	}

	c.String(http.StatusAccepted, _orderAccepted)
}

func (oh *OrderHandler) ListOrders(c *gin.Context) {
	dbItems, err := oh.stg.ListOrders(c.Request.Context(), auth.GetUserID(c))
	if err != nil {
		if errors.Is(storage.ErrNoOrders, err) {
			_http.CreateStatusResponse(c, http.StatusNoContent)
			return
		}

		_http.CreateStatusResponse(c, http.StatusInternalServerError)
		return
	}

	respItems := make([]OrderItemsResponse, 0, len(dbItems))
	for _, dbItem := range dbItems {
		respItems = append(respItems, OrderItemsResponse{
			Number:     dbItem.Number,
			Status:     dbItem.Status,
			Accrual:    dbItem.Accrual,
			UploadedAt: dbItem.UploadedAt,
		})
	}

	c.JSON(http.StatusOK, respItems)
}
