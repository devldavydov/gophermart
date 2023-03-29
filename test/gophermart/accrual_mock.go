package gophermart

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type AccrualMockResp map[string]map[string]interface{}

type AccrualMock struct {
	mu            sync.RWMutex
	respMap       AccrualMockResp
	listenAddress string
}

func NewAccrualMock(listenAddress string) *AccrualMock {
	return &AccrualMock{
		respMap:       make(AccrualMockResp),
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

func (am *AccrualMock) SetRespMap(respMap AccrualMockResp) {
	data := make(AccrualMockResp, len(respMap))
	for k, v := range respMap {
		data[k] = v
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	am.respMap = respMap
}

func (am *AccrualMock) handlerOrder(c *gin.Context) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	resp, ok := am.respMap[c.Param("number")]
	if !ok {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}
