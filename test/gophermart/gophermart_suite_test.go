package gophermart

import (
	"fmt"
	"net/http/cookiejar"
	"os"
	"testing"

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
	accrualSrvListenAddr string
}

func (gs *GophermartSuite) SetupSuite() {
	jar, _ := cookiejar.New(nil)
	gs.httpClient = resty.New().
		SetHostURL(os.Getenv(_envGophermartSrvAddr)).
		SetCookieJar(jar)

	gs.accrualSrvListenAddr = os.Getenv(_envAccrualSrvListenAddr)
}

func (gs *GophermartSuite) TearDownSuite() {
	fmt.Println("stop")
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
