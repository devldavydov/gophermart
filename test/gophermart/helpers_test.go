package gophermart

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func (gs *GophermartSuite) userRegister(userLogin, userPassword string, statusCode int) *resty.Response {
	resp, err := gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/register")
	require.NoError(gs.T(), err)
	require.Equal(gs.T(), statusCode, resp.StatusCode())
	return resp
}

func (gs *GophermartSuite) userLogin(userLogin, userPassword string, statusCode int) {
	resp, err := gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/login")
	require.NoError(gs.T(), err)
	require.Equal(gs.T(), statusCode, resp.StatusCode())
}

func (gs *GophermartSuite) userLogout() {
	resp, err := gs.httpClient.R().Post("/api/user/logout")
	require.NoError(gs.T(), err)
	require.Equal(gs.T(), http.StatusOK, resp.StatusCode())
}

func (gs *GophermartSuite) userAddOrder(orderNum string, statusCode int) {
	resp, err := gs.httpClient.R().
		SetBody(orderNum).
		SetHeader("Content-Type", "text/plain").
		Post("/api/user/orders")
	require.NoError(gs.T(), err)
	require.Equal(gs.T(), statusCode, resp.StatusCode())
}

func (gs *GophermartSuite) userGetOrders(statusCode int) orderList {
	var lst orderList
	resp, err := gs.httpClient.R().
		SetResult(&lst).
		Get("/api/user/orders")
	require.NoError(gs.T(), err)
	require.Equal(gs.T(), statusCode, resp.StatusCode())
	return lst
}

func (gs *GophermartSuite) userGetBalance(statusCode int) userBalance {
	var blnc userBalance
	resp, err := gs.httpClient.R().
		SetResult(&blnc).
		Get("/api/user/balance")
	require.NoError(gs.T(), err)
	require.Equal(gs.T(), statusCode, resp.StatusCode())
	return blnc
}

func (gs *GophermartSuite) userGetBalanceWithdrawals(statusCode int) userBalanceWithdrawals {
	var wthdrwls userBalanceWithdrawals
	resp, err := gs.httpClient.R().
		SetResult(&wthdrwls).
		Get("/api/user/withdrawals")
	require.NoError(gs.T(), err)
	require.Equal(gs.T(), statusCode, resp.StatusCode())
	return wthdrwls
}

func (gs *GophermartSuite) userBalanceWithdraw(blnc userBalanceWithdraw, statusCode int) {
	resp, err := gs.httpClient.R().
		SetBody(blnc).
		Post("/api/user/balance/withdraw")
	require.NoError(gs.T(), err)
	require.Equal(gs.T(), statusCode, resp.StatusCode())
}

func (gs *GophermartSuite) startAccrualMock(wg *sync.WaitGroup, respMap AccrualMockResp) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		NewAccrualMock(respMap, gs.accrualSrvListenAddr).Start(ctx)
	}()

	return cancel
}
