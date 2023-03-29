package gophermart

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/devldavydov/gophermart/pkg/gophermart"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func (gs *GophermartSuite) TestBalanceApi() {
	require.NoError(gs.T(), gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))

	gs.Run("get user balance", func() {
		blnc, err := gs.gCli.GetBalance()
		gs.NoError(err)
		gs.Equal(float64(0), blnc.Current)
		gs.Equal(float64(0), blnc.Withdrawn)
	})

	gs.Run("get user withdrawals", func() {
		_, err := gs.gCli.GetBalanceWithdrawals()
		gs.ErrorIs(gophermart.ErrNoBalanceWithdrawals, err)
	})

	gs.Run("balance withdraw wrong format", func() {
		gs.ErrorIs(gophermart.ErrBalanceWithdrawWrongFormat, gs.gCli.BalanceWithdraw("123", 123))
	})

	gs.Run("balance withdraw not enough balance", func() {
		gs.ErrorIs(gophermart.ErrBalanceWithdrawPaymentRequired, gs.gCli.BalanceWithdraw(goluhn.Generate(10), 123))
	})
}
