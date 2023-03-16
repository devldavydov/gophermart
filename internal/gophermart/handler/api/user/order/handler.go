package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (oh *OrderHandler) AddOrder(c *gin.Context) {
	c.String(http.StatusOK, "AddOrder\n")
}

func (oh *OrderHandler) ListOrders(c *gin.Context) {
	c.String(http.StatusOK, "ListOrders\n")
}
