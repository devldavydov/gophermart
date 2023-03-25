package gophermart

import (
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func (gs *GophermartSuite) TestBalanceApi() {
	userLogin, userPassword := uuid.NewString(), uuid.NewString()

	// Register user
	resp, err := gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/register")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())

	// Get user balance
	var blnc userBalance
	resp, err = gs.httpClient.R().
		SetResult(&blnc).
		Get("/api/user/balance")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())
	assert.Equal(gs.T(), float64(0), blnc.Current)
	assert.Equal(gs.T(), float64(0), blnc.Withdrawn)

	// Get user withdrawls
	resp, err = gs.httpClient.R().
		Get("/api/user/withdrawals")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusNoContent, resp.StatusCode())

	// Balance withdraw - wrong request
	resp, err = gs.httpClient.R().
		SetBody("foobar").
		Post("/api/user/balance/withdraw")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusBadRequest, resp.StatusCode())

	// Balance withdraw - wrong format
	resp, err = gs.httpClient.R().
		SetBody(userBalanceWithdraw{Order: "123", Sum: 123}).
		Post("/api/user/balance/withdraw")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusUnprocessableEntity, resp.StatusCode())

	// Balance withdraw - not enough balance
	resp, err = gs.httpClient.R().
		SetBody(userBalanceWithdraw{Order: goluhn.Generate(10), Sum: 123}).
		Post("/api/user/balance/withdraw")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusPaymentRequired, resp.StatusCode())
}
