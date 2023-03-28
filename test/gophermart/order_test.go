package gophermart

import (
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func (gs *GophermartSuite) TestOrderApi() {
	userLogin, userPassword := uuid.NewString(), uuid.NewString()
	orderNums := []string{goluhn.Generate(10), goluhn.Generate(10)}

	resp, err := gs.userRegister(userLogin, userPassword)
	require.NoError(gs.T(), err)
	require.Equal(gs.T(), http.StatusOK, resp.StatusCode())

	gs.Run("add order wrong request", func() {
		resp, err = gs.httpClient.R().
			SetBody(userAuth{Login: userLogin, Password: userPassword}).
			Post("/api/user/orders")
		gs.NoError(err)
		gs.Equal(http.StatusBadRequest, resp.StatusCode())
	})

	gs.Run("add order wrong format", func() {
		resp, err = gs.httpClient.R().
			SetBody("123").
			SetHeader("Content-Type", "text/plain").
			Post("/api/user/orders")
		gs.NoError(err)
		gs.Equal(http.StatusUnprocessableEntity, resp.StatusCode())
	})

	gs.Run("get empty orders list", func() {
		resp, err = gs.httpClient.R().
			Get("/api/user/orders")
		gs.NoError(err)
		gs.Equal(http.StatusNoContent, resp.StatusCode())
	})

	gs.Run("add orders", func() {
		for _, orderNum := range orderNums {
			resp, err = gs.userAddOrder(orderNum)
			gs.NoError(err)
			gs.Equal(http.StatusAccepted, resp.StatusCode())
		}

	})

	gs.Run("add same order", func() {
		resp, err = gs.userAddOrder(orderNums[0])
		gs.NoError(err)
		gs.Equal(http.StatusOK, resp.StatusCode())
	})

	gs.Run("get orders list", func() {
		var lst orderList
		resp, err = gs.httpClient.R().
			SetResult(&lst).
			Get("/api/user/orders")
		gs.NoError(err)
		gs.Equal(http.StatusOK, resp.StatusCode())
		gs.Equal(2, len(lst))

		gs.Equal(orderNums[0], lst[0].Number)
		gs.True(lst[0].Status == "NEW" || lst[1].Status == "PROCESSING")
		gs.Equal(orderNums[1], lst[1].Number)
		gs.True(lst[1].Status == "NEW" || lst[1].Status == "PROCESSING")
	})

	gs.Run("add existing order from another user", func() {
		resp, err = gs.httpClient.R().
			SetBody(userAuth{Login: uuid.NewString(), Password: uuid.NewString()}).
			Post("/api/user/register")
		gs.NoError(err)
		gs.Equal(http.StatusOK, resp.StatusCode())

		resp, err = gs.userAddOrder(orderNums[0])
		gs.NoError(err)
		gs.Equal(http.StatusConflict, resp.StatusCode())
	})
}
