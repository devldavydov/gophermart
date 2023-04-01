package login

import (
	"errors"
	"net/http"

	_http "github.com/devldavydov/gophermart/internal/common/http"
	"github.com/devldavydov/gophermart/internal/gophermart/auth"
	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	_userAlreadyExists = "User already exists"
	_userWrongAuth     = "User wrong login/password"
)

type LoginReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginHandler struct {
	stg    storage.Storage
	logger *logrus.Logger
}

func NewLoginHandler(stg storage.Storage, logger *logrus.Logger) *LoginHandler {
	return &LoginHandler{stg: stg, logger: logger}
}

func (lh *LoginHandler) Register(c *gin.Context) {
	var req LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		_http.CreateStatusResponse(c, http.StatusBadRequest)
		return
	}

	pwdHash, err := hashPassword(req.Password)
	if err != nil {
		lh.logger.Errorf("failed to get password hash: %v", err)
		_http.CreateStatusResponse(c, http.StatusInternalServerError)
		return
	}

	userID, err := lh.stg.CreateUser(c.Request.Context(), req.Login, pwdHash)
	if err != nil {
		if errors.Is(storage.ErrUserAlreadyExists, err) {
			c.String(http.StatusConflict, _userAlreadyExists)
			return
		}

		_http.CreateStatusResponse(c, http.StatusInternalServerError)
		return
	}

	if err = auth.SetUserID(c, userID); err != nil {
		_http.CreateStatusResponse(c, http.StatusInternalServerError)
		return
	}

	_http.CreateStatusResponse(c, http.StatusOK)
}

func (lh *LoginHandler) Login(c *gin.Context) {
	var req LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		_http.CreateStatusResponse(c, http.StatusBadRequest)
		return
	}

	userID, pwdHash, err := lh.stg.FindUser(c.Request.Context(), req.Login)
	if err != nil {
		if errors.Is(storage.ErrUserNotFound, err) {
			c.String(http.StatusUnauthorized, _userWrongAuth)
			return
		}

		_http.CreateStatusResponse(c, http.StatusInternalServerError)
		return
	}

	err = checkPassword(req.Password, pwdHash)
	if err != nil {
		c.String(http.StatusUnauthorized, _userWrongAuth)
		return
	}

	if err = auth.SetUserID(c, userID); err != nil {
		_http.CreateStatusResponse(c, http.StatusInternalServerError)
		return
	}

	_http.CreateStatusResponse(c, http.StatusOK)
}

func (lh *LoginHandler) Logout(c *gin.Context) {
	if err := auth.DelUserID(c); err != nil {
		_http.CreateStatusResponse(c, http.StatusInternalServerError)
		return
	}
	_http.CreateStatusResponse(c, http.StatusOK)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	return string(bytes), err
}

func checkPassword(password, pwdHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(password))
}
