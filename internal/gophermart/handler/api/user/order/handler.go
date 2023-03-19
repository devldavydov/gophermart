package order

import (
	"net/http"

	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	stg    storage.Storage
	logger *logrus.Logger
}

func NewOrderHandler(stg storage.Storage, logger *logrus.Logger) *OrderHandler {
	return &OrderHandler{stg: stg, logger: logger}
}

func (oh *OrderHandler) AddOrder(c *gin.Context) {
	c.String(http.StatusOK, "AddOrder\n")
}

func (oh *OrderHandler) ListOrders(c *gin.Context) {
	c.String(http.StatusOK, "ListOrders\n")
}
