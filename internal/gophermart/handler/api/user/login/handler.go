package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct{}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (lh *LoginHandler) Register(c *gin.Context) {
	c.String(http.StatusOK, "Register\n")
}

func (lh *LoginHandler) Login(c *gin.Context) {
	c.String(http.StatusOK, "Login\n")
}

func (lh *LoginHandler) Logout(c *gin.Context) {
	c.String(http.StatusOK, "Logout\n")
}
