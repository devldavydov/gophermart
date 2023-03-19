package order

import (
	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Init(group *gin.RouterGroup, stg storage.Storage, logger *logrus.Logger) {
	orderHandler := NewOrderHandler(stg, logger)
	group.POST("/orders", orderHandler.AddOrder)
	group.GET("/orders", orderHandler.ListOrders)
}
