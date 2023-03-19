package order

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	_http "github.com/devldavydov/gophermart/internal/common/http"
	"github.com/devldavydov/gophermart/internal/gophermart/auth"
	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/theplant/luhn"
)

const (
	_orderAlreadyExists         = "Order already exists"
	_orderAlreadyExistsFromUser = "Order already exists from user"
	_orderAccepted              = "Order accepted"
)

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

	orderNumBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_http.CreateStatusResponse(c, http.StatusBadRequest)
		return
	}

	orderNum := string(orderNumBytes)
	if !checkNumLuhn(orderNum) {
		_http.CreateStatusResponse(c, http.StatusUnprocessableEntity)
		return
	}

	err = oh.stg.AddOrder(auth.GetUserId(c), orderNum)
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
	c.String(http.StatusOK, "ListOrders\n")
}

func checkNumLuhn(orderNum string) bool {
	order, err := strconv.Atoi(orderNum)
	if err != nil {
		return false
	}
	return luhn.Valid(order)
}
