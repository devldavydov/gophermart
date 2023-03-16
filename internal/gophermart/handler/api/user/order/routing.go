package order

import "github.com/gin-gonic/gin"

func Init(group *gin.RouterGroup) {
	orderHandler := NewOrderHandler()
	group.POST("/orders", orderHandler.AddOrder)
	group.GET("/orders", orderHandler.ListOrders)
}
