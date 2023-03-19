package login

import (
	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Init(group *gin.RouterGroup, stg storage.Storage, logger *logrus.Logger) {
	loginHandler := NewLoginHandler(stg, logger)
	group.POST("/register", loginHandler.Register)
	group.POST("/login", loginHandler.Login)
	group.POST("/logout", loginHandler.Logout)
}
