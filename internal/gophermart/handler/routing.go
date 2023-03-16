package handler

import (
	"github.com/devldavydov/gophermart/internal/gophermart/handler/api/user/balance"
	"github.com/devldavydov/gophermart/internal/gophermart/handler/api/user/login"
	"github.com/devldavydov/gophermart/internal/gophermart/handler/api/user/order"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	apiUser := router.Group("/api/user")

	balance.Init(apiUser)
	login.Init(apiUser)
	order.Init(apiUser)
}
