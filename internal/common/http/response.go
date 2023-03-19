package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateStatusResponse(c *gin.Context, statusCode int) {
	c.String(statusCode, http.StatusText(statusCode))
}
