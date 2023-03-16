package balance

import "github.com/gin-gonic/gin"

func Init(group *gin.RouterGroup) {
	balanceHandler := NewBalanceHandler()
	group.GET("/balance", balanceHandler.GetBalance)
	group.POST("/balance/withdraw", balanceHandler.BalanceWithdraw)
	group.GET("/withdrawals", balanceHandler.ListWithdrawals)
}
