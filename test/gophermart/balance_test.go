package gophermart

import (
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func (gs *GophermartSuite) TestBalanceApi() {
	userLogin, userPassword := uuid.NewString(), uuid.NewString()

	resp, err := gs.userRegister(userLogin, userPassword)
	require.NoError(gs.T(), err)
	require.Equal(gs.T(), http.StatusOK, resp.StatusCode())

	gs.Run("get user balance", func() {
		var blnc userBalance
		resp, err = gs.httpClient.R().
			SetResult(&blnc).
			Get("/api/user/balance")
		gs.NoError(err)
		gs.Equal(http.StatusOK, resp.StatusCode())
		gs.Equal(float64(0), blnc.Current)
		gs.Equal(float64(0), blnc.Withdrawn)
	})

	gs.Run("get user withdrawals", func() {
		resp, err = gs.httpClient.R().
			Get("/api/user/withdrawals")
		gs.NoError(err)
		gs.Equal(http.StatusNoContent, resp.StatusCode())
	})

	gs.Run("balance withdraw wrong request", func() {
		resp, err = gs.httpClient.R().
			SetBody("foobar").
			Post("/api/user/balance/withdraw")
		gs.NoError(err)
		gs.Equal(http.StatusBadRequest, resp.StatusCode())
	})

	gs.Run("balance withdraw wrong format", func() {
		resp, err = gs.httpClient.R().
			SetBody(userBalanceWithdraw{Order: "123", Sum: 123}).
			Post("/api/user/balance/withdraw")
		gs.NoError(err)
		gs.Equal(http.StatusUnprocessableEntity, resp.StatusCode())
	})

	gs.Run("balance withdraw not enough balance", func() {
		resp, err = gs.httpClient.R().
			SetBody(userBalanceWithdraw{Order: goluhn.Generate(10), Sum: 123}).
			Post("/api/user/balance/withdraw")
		gs.NoError(err)
		gs.Equal(http.StatusPaymentRequired, resp.StatusCode())
	})
}
