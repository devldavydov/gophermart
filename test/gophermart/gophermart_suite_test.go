package gophermart

import (
	"context"
	"net/http/cookiejar"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/devldavydov/gophermart/pkg/gophermart"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

const (
	_envGophermartSrvAddr    = "TEST_GOPHERMART_SRV_ADDR"
	_envAccrualSrvListenAddr = "TEST_ACCRUAL_SRV_LISTEN_ADDR"
)

type GophermartSuite struct {
	suite.Suite

	httpClient           *resty.Client
	gCli                 *gophermart.Client
	accrualSrvListenAddr string
	waitTimeout          time.Duration
	waitTick             time.Duration
	wg                   sync.WaitGroup
	accrualStop          context.CancelFunc
	accrualMock          *AccrualMock
}

func (gs *GophermartSuite) SetupSuite() {
	jar, _ := cookiejar.New(nil)
	gs.httpClient = resty.New().
		SetBaseURL(os.Getenv(_envGophermartSrvAddr)).
		SetCookieJar(jar)

	gs.gCli = gophermart.NewClient(os.Getenv(_envGophermartSrvAddr))

	ctx, stop := context.WithCancel(context.Background())
	gs.accrualMock = NewAccrualMock(os.Getenv(_envAccrualSrvListenAddr))

	gs.wg.Add(1)
	go func() {
		defer gs.wg.Done()
		gs.accrualMock.Start(ctx)
	}()

	gs.accrualStop = stop
	gs.waitTimeout = 1 * time.Minute
	gs.waitTick = 1 * time.Second
}

func (gs *GophermartSuite) TearDownSuite() {
	gs.accrualStop()
	gs.wg.Wait()
}

func TestGophermartIntegration(t *testing.T) {
	_, ok1 := os.LookupEnv(_envGophermartSrvAddr)
	_, ok2 := os.LookupEnv(_envAccrualSrvListenAddr)
	if !(ok1 && ok2) {
		t.Skip("Test environment no set")
		return
	}

	suite.Run(t, new(GophermartSuite))
}
