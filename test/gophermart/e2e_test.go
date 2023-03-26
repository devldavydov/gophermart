package gophermart

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func (gs *GophermartSuite) TestE2E() {
	userLogin, userPassword := uuid.NewString(), uuid.NewString()

	orderNum1 := goluhn.Generate(10)
	orderNum2 := goluhn.Generate(10)
	orderNum3 := goluhn.Generate(10)
	orderNum4 := goluhn.Generate(10)
	orderNum5 := goluhn.Generate(10)

	// Prepare accrual response map
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	respMap := AccrualMockResp{
		orderNum1: {
			"order":  orderNum1,
			"status": "REGISTERED",
		},
		orderNum2: {
			"order":  orderNum2,
			"status": "INVALID",
		},
		orderNum3: {
			"order":  orderNum3,
			"status": "PROCESSING",
		},
		orderNum4: {
			"order":   orderNum4,
			"status":  "PROCESSED",
			"accrual": 100.5,
		},
		orderNum5: {
			"order":   orderNum5,
			"status":  "PROCESSED",
			"accrual": 100.55,
		},
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		NewAccrualMock(respMap, gs.accrualSrvListenAddr).Start(ctx)
	}()

	// Register user
	resp, err := gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/register")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())

	// Get balance
	var blnc userBalance
	resp, err = gs.httpClient.R().
		SetResult(&blnc).
		Get("/api/user/balance")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())
	assert.Equal(gs.T(), float64(0), blnc.Current)
	assert.Equal(gs.T(), float64(0), blnc.Withdrawn)

	// Add orders
	for _, order := range []string{orderNum1, orderNum2, orderNum3, orderNum4, orderNum5} {
		resp, err = gs.httpClient.R().
			SetBody(order).
			SetHeader("Content-Type", "text/plain").
			Post("/api/user/orders")
		assert.NoError(gs.T(), err)
		assert.Equal(gs.T(), http.StatusAccepted, resp.StatusCode())
	}

	// Wait for balance computation
	waitDuration := 2 * time.Minute
	startTime := time.Now()
	success := false
	for time.Since(startTime) < waitDuration {
		resp, err = gs.httpClient.R().
			SetResult(&blnc).
			Get("/api/user/balance")

		if err == nil &&
			resp.StatusCode() == http.StatusOK &&
			float64(201.05) == blnc.Current {
			success = true
			break
		}

		time.Sleep(1 * time.Second)
	}
	assert.True(gs.T(), success)

	// Get orders
	var lst orderList
	resp, err = gs.httpClient.R().
		SetResult(&lst).
		Get("/api/user/orders")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())
	assert.Equal(gs.T(), 5, len(lst))

	for _, item := range lst {
		switch item.Number {
		case orderNum2:
			assert.Equal(gs.T(), "INVALID", item.Status)
		case orderNum4:
			assert.Equal(gs.T(), "PROCESSED", item.Status)
			assert.Equal(gs.T(), 100.5, *item.Accrual)
		case orderNum5:
			assert.Equal(gs.T(), "PROCESSED", item.Status)
			assert.Equal(gs.T(), 100.55, *item.Accrual)
		}
	}

	// Balance withdraw
	orderWithdraw := goluhn.Generate(10)
	resp, err = gs.httpClient.R().
		SetBody(userBalanceWithdraw{Order: orderWithdraw, Sum: 101.04}).
		Post("/api/user/balance/withdraw")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())

	// Check new balance
	resp, err = gs.httpClient.R().
		SetResult(&blnc).
		Get("/api/user/balance")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())
	assert.Equal(gs.T(), float64(100.01), blnc.Current)
	assert.Equal(gs.T(), float64(101.04), blnc.Withdrawn)

	// Get withdrawals
	var wthdrwls userBalanceWithdrawals
	resp, err = gs.httpClient.R().
		SetResult(&wthdrwls).
		Get("/api/user/withdrawals")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())
	assert.Equal(gs.T(), 1, len(wthdrwls))

	assert.Equal(gs.T(), orderWithdraw, wthdrwls[0].Order)
	assert.Equal(gs.T(), 101.04, wthdrwls[0].Sum)

	cancel()
	wg.Wait()
}
