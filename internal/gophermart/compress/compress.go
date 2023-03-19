package compress

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.Use(gzip.Gzip(gzip.DefaultCompression))
}
