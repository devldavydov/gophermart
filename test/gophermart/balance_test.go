package gophermart

import (
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/google/uuid"
)

func (gs *GophermartSuite) TestBalanceApi() {
	userLogin, userPassword := uuid.NewString(), uuid.NewString()

	gs.userRegister(userLogin, userPassword, http.StatusOK)

	gs.Run("get user balance", func() {
		blnc := gs.userGetBalance(http.StatusOK)
		gs.Equal(float64(0), blnc.Current)
		gs.Equal(float64(0), blnc.Withdrawn)
	})

	gs.Run("get user withdrawals", func() {
		gs.userGetBalanceWithdrawals(http.StatusNoContent)
	})

	gs.Run("balance withdraw wrong request", func() {
		resp, err := gs.httpClient.R().
			SetBody("foobar").
			Post("/api/user/balance/withdraw")
		gs.NoError(err)
		gs.Equal(http.StatusBadRequest, resp.StatusCode())
	})

	gs.Run("balance withdraw wrong format", func() {
		gs.userBalanceWithdraw(userBalanceWithdraw{Order: "123", Sum: 123}, http.StatusUnprocessableEntity)
	})

	gs.Run("balance withdraw not enough balance", func() {
		gs.userBalanceWithdraw(userBalanceWithdraw{Order: goluhn.Generate(10), Sum: 123}, http.StatusPaymentRequired)
	})
}
