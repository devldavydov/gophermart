package gophermart

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func (gs *GophermartSuite) TestE2EOrderInvalid() {
	require.NoError(gs.T(), gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))
	orderNum := goluhn.Generate(10)
	gs.accrualMock.SetRespMap(AccrualMockResp{orderNum: {"order": orderNum, "status": "INVALID"}})

	gs.Run("add order", func() {
		gs.NoError(gs.gCli.AddOrder(orderNum))
	})

	gs.Run("wait order status changed", func() {
		gs.Eventually(func() bool {
			lst, err := gs.gCli.GetOrders()
			if err != nil {
				return false
			}

			return len(lst) == 1 &&
				lst[0].Status == "INVALID" &&
				lst[0].Number == orderNum &&
				lst[0].Accrual == nil
		}, gs.waitTimeout, gs.waitTick)
	})
}

func (gs *GophermartSuite) TestE2EOrderProcessed() {
	require.NoError(gs.T(), gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))
	orderNum := goluhn.Generate(10)
	gs.accrualMock.SetRespMap(AccrualMockResp{orderNum: {"order": orderNum, "status": "PROCESSED", "accrual": 100}})

	gs.Run("add order", func() {
		gs.NoError(gs.gCli.AddOrder(orderNum))
	})

	gs.Run("wait order status changed", func() {
		gs.Eventually(func() bool {
			lst, err := gs.gCli.GetOrders()
			if err != nil {
				return false
			}

			return len(lst) == 1 &&
				lst[0].Status == "PROCESSED" &&
				lst[0].Number == orderNum &&
				lst[0].Accrual != nil &&
				*lst[0].Accrual == 100
		}, gs.waitTimeout, gs.waitTick)
	})
}

func (gs *GophermartSuite) TestE2ECheckBalance() {
	require.NoError(gs.T(), gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))
	orderNum := goluhn.Generate(10)
	gs.accrualMock.SetRespMap(AccrualMockResp{orderNum: {"order": orderNum, "status": "PROCESSED", "accrual": 100}})

	gs.Run("check initial balance", func() {
		blnc, err := gs.gCli.GetBalance()
		gs.NoError(err)
		gs.Equal(float64(0), blnc.Current)
		gs.Equal(float64(0), blnc.Withdrawn)
	})

	gs.Run("add order", func() {
		gs.NoError(gs.gCli.AddOrder(orderNum))
	})

	gs.Run("wait balance changed", func() {
		gs.Eventually(func() bool {
			blnc, err := gs.gCli.GetBalance()
			if err != nil {
				return false
			}

			return blnc.Current == float64(100) &&
				blnc.Withdrawn == float64(0)
		}, gs.waitTimeout, gs.waitTick)
	})
}

func (gs *GophermartSuite) TestE2ECheckWithdrawals() {
	require.NoError(gs.T(), gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))
	orderNum1, orderNum2, orderWdraw := goluhn.Generate(10), goluhn.Generate(10), goluhn.Generate(10)
	gs.accrualMock.SetRespMap(AccrualMockResp{
		orderNum1: {"order": orderNum1, "status": "PROCESSED", "accrual": 100},
		orderNum2: {"order": orderNum2, "status": "PROCESSED", "accrual": 100},
	})

	gs.Run("add orders", func() {
		gs.NoError(gs.gCli.AddOrder(orderNum1))
		gs.NoError(gs.gCli.AddOrder(orderNum2))
	})

	gs.Run("wait balance changed", func() {
		gs.Eventually(func() bool {
			blnc, err := gs.gCli.GetBalance()
			if err != nil {
				return false
			}

			return blnc.Current == float64(200) &&
				blnc.Withdrawn == float64(0)
		}, gs.waitTimeout, gs.waitTick)
	})

	gs.Run("balance withdraw", func() {
		gs.gCli.BalanceWithdraw(orderWdraw, 100.1)
	})

	gs.Run("check changed balance", func() {
		blnc, err := gs.gCli.GetBalance()
		gs.NoError(err)
		gs.Equal(float64(99.9), blnc.Current)
		gs.Equal(float64(100.1), blnc.Withdrawn)
	})

	gs.Run("check withdrawals", func() {
		lst, err := gs.gCli.GetBalanceWithdrawals()
		gs.NoError(err)
		gs.Equal(1, len(lst))

		gs.Equal(orderWdraw, lst[0].Order)
		gs.Equal(float64(100.1), lst[0].Sum)
	})
}

func (gs *GophermartSuite) TestE2ECheckBalanceWithFailedOrders() {
	require.NoError(gs.T(), gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))
	orders := []string{goluhn.Generate(10), goluhn.Generate(10), goluhn.Generate(10), goluhn.Generate(10), goluhn.Generate(10)}
	gs.accrualMock.SetRespMap(AccrualMockResp{
		orders[0]: {"order": orders[0], "status": "PROCESSED", "accrual": 100},
		orders[1]: {"order": orders[1], "status": "INVALID"},
		orders[2]: {"order": orders[2], "status": "PROCESSING"},
		orders[3]: {"order": orders[3], "status": "REGISTERED"},
	})

	gs.Run("add orders", func() {
		for _, order := range orders {
			gs.NoError(gs.gCli.AddOrder(order))
		}
	})

	gs.Run("wait balance changed", func() {
		gs.Eventually(func() bool {
			blnc, err := gs.gCli.GetBalance()
			if err != nil {
				return false
			}

			return blnc.Current == float64(100) &&
				blnc.Withdrawn == float64(0)
		}, gs.waitTimeout, gs.waitTick)
	})
}
