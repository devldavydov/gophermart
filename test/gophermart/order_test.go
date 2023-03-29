package gophermart

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/devldavydov/gophermart/pkg/gophermart"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func (gs *GophermartSuite) TestOrderApi() {
	orderNums := []string{goluhn.Generate(10), goluhn.Generate(10)}

	require.NoError(gs.T(), gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))

	gs.Run("add order wrong format", func() {
		gs.ErrorIs(gophermart.ErrOrderWrongFormat, gs.gCli.AddOrder("123"))
	})

	gs.Run("get empty orders list", func() {
		_, err := gs.gCli.GetOrders()
		gs.ErrorIs(gophermart.ErrNoOrders, err)
	})

	gs.Run("add orders", func() {
		for _, orderNum := range orderNums {
			gs.NoError(gs.gCli.AddOrder(orderNum))
		}
	})

	gs.Run("add same order", func() {
		gs.ErrorIs(gophermart.ErrOrderAlreadyAccepted, gs.gCli.AddOrder(orderNums[0]))
	})

	gs.Run("get orders list", func() {
		lst, err := gs.gCli.GetOrders()
		gs.NoError(err)
		gs.Equal(2, len(lst))

		gs.Equal(orderNums[0], lst[0].Number)
		gs.True(lst[0].Status == "NEW" || lst[1].Status == "PROCESSING")
		gs.Equal(orderNums[1], lst[1].Number)
		gs.True(lst[1].Status == "NEW" || lst[1].Status == "PROCESSING")
	})

	gs.Run("add existing order from another user", func() {
		gs.NoError(gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))
		gs.ErrorIs(gophermart.ErrOrderAlreadyExists, gs.gCli.AddOrder(orderNums[0]))
	})
}
