package handler

import (
	"github.com/devldavydov/gophermart/internal/gophermart/handler/api/user/balance"
	"github.com/devldavydov/gophermart/internal/gophermart/handler/api/user/login"
	"github.com/devldavydov/gophermart/internal/gophermart/handler/api/user/order"
	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Init(router *gin.Engine, stg storage.Storage, logger *logrus.Logger) {
	apiUser := router.Group("/api/user")

	balance.Init(apiUser, stg, logger)
	login.Init(apiUser, stg, logger)
	order.Init(apiUser, stg, logger)
}
