package gophermart

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AccrualMockResp map[string]map[string]interface{}

type AccrualMock struct {
	respMap       AccrualMockResp
	listenAddress string
}

func NewAccrualMock(respMap AccrualMockResp, listenAddress string) *AccrualMock {
	return &AccrualMock{
		respMap:       respMap,
		listenAddress: listenAddress,
	}
}

func (am *AccrualMock) Start(ctx context.Context) {
	router := gin.Default()
	router.GET("/api/orders/:number", am.handlerOrder)

	httpServer := &http.Server{
		Addr:    am.listenAddress,
		Handler: router,
	}

	go httpServer.ListenAndServe()
	<-ctx.Done()

	hCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	httpServer.Shutdown(hCtx)
}

func (am *AccrualMock) handlerOrder(c *gin.Context) {
	number := c.Param("number")

	resp, ok := am.respMap[number]
	if !ok {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}
