package auth

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var ErrUserSessionNotFound = errors.New("user session not found")

const (
	_userSession    = "UserSession"
	_userSessionKey = "UserId"
)

func Init(router *gin.Engine, secret string) {
	router.Use(sessions.Sessions(_userSession, cookie.NewStore([]byte(secret))))
}

func AuthRequired(c *gin.Context) {
	if sessions.Default(c).Get(_userSessionKey) == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

func SetUserID(c *gin.Context, userID int) error {
	session := sessions.Default(c)
	session.Set(_userSessionKey, userID)
	return session.Save()
}

func GetUserID(c *gin.Context) int {
	val := sessions.Default(c).Get(_userSessionKey)
	userID, _ := val.(int)
	return userID
}

func DelUserID(c *gin.Context) error {
	session := sessions.Default(c)
	val := session.Get(_userSessionKey)
	if val == nil {
		return nil
	}
	session.Delete(_userSessionKey)
	return session.Save()
}
