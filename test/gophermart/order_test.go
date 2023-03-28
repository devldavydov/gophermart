package gophermart

import (
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/google/uuid"
)

func (gs *GophermartSuite) TestOrderApi() {
	userLogin, userPassword := uuid.NewString(), uuid.NewString()
	orderNums := []string{goluhn.Generate(10), goluhn.Generate(10)}

	gs.userRegister(userLogin, userPassword, http.StatusOK)

	gs.Run("add order wrong request", func() {
		resp, err := gs.httpClient.R().
			SetBody(userAuth{Login: userLogin, Password: userPassword}).
			Post("/api/user/orders")
		gs.NoError(err)
		gs.Equal(http.StatusBadRequest, resp.StatusCode())
	})

	gs.Run("add order wrong format", func() {
		resp, err := gs.httpClient.R().
			SetBody("123").
			SetHeader("Content-Type", "text/plain").
			Post("/api/user/orders")
		gs.NoError(err)
		gs.Equal(http.StatusUnprocessableEntity, resp.StatusCode())
	})

	gs.Run("get empty orders list", func() {
		gs.userGetOrders(http.StatusNoContent)
	})

	gs.Run("add orders", func() {
		for _, orderNum := range orderNums {
			gs.userAddOrder(orderNum, http.StatusAccepted)
		}
	})

	gs.Run("add same order", func() {
		gs.userAddOrder(orderNums[0], http.StatusOK)
	})

	gs.Run("get orders list", func() {
		lst := gs.userGetOrders(http.StatusOK)
		gs.Equal(2, len(lst))

		gs.Equal(orderNums[0], lst[0].Number)
		gs.True(lst[0].Status == "NEW" || lst[1].Status == "PROCESSING")
		gs.Equal(orderNums[1], lst[1].Number)
		gs.True(lst[1].Status == "NEW" || lst[1].Status == "PROCESSING")
	})

	gs.Run("add existing order from another user", func() {
		gs.userRegister(uuid.NewString(), uuid.NewString(), http.StatusOK)
		gs.userAddOrder(orderNums[0], http.StatusConflict)
	})
}
