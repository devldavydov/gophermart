package login

import "github.com/gin-gonic/gin"

func Init(group *gin.RouterGroup) {
	loginHandler := NewLoginHandler()
	group.POST("/register", loginHandler.Register)
	group.POST("/login", loginHandler.Login)
	group.POST("/logout", loginHandler.Logout)
}
