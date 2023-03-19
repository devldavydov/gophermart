package balance

import (
	"github.com/devldavydov/gophermart/internal/gophermart/auth"
	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Init(group *gin.RouterGroup, stg storage.Storage, logger *logrus.Logger) {
	balanceHandler := NewBalanceHandler(stg, logger)
	group.GET("/balance", auth.AuthRequired, balanceHandler.GetBalance)
	group.POST("/balance/withdraw", balanceHandler.BalanceWithdraw)
	group.GET("/withdrawals", balanceHandler.ListWithdrawals)
}
