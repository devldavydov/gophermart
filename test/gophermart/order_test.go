package gophermart

import (
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func (gs *GophermartSuite) TestOrderApi() {
	userLogin, userPassword := uuid.NewString(), uuid.NewString()

	// Register user
	resp, err := gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/register")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())

	// Add order wrong request
	resp, err = gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/orders")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusBadRequest, resp.StatusCode())

	// Add order wrong format
	resp, err = gs.httpClient.R().
		SetBody("123").
		SetHeader("Content-Type", "text/plain").
		Post("/api/user/orders")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusUnprocessableEntity, resp.StatusCode())

	orderNum1, orderNum2 := goluhn.Generate(10), goluhn.Generate(10)

	// Add order
	resp, err = gs.httpClient.R().
		SetBody(orderNum1).
		SetHeader("Content-Type", "text/plain").
		Post("/api/user/orders")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusAccepted, resp.StatusCode())

	// Add same order
	resp, err = gs.httpClient.R().
		SetBody(orderNum1).
		SetHeader("Content-Type", "text/plain").
		Post("/api/user/orders")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())

	// Add another order
	resp, err = gs.httpClient.R().
		SetBody(orderNum2).
		SetHeader("Content-Type", "text/plain").
		Post("/api/user/orders")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusAccepted, resp.StatusCode())

	// Get orders list
	var lst orderList
	resp, err = gs.httpClient.R().
		SetResult(&lst).
		Get("/api/user/orders")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())
	assert.Equal(gs.T(), 2, len(lst))

	assert.Equal(gs.T(), orderNum1, lst[0].Number)
	assert.True(gs.T(), lst[0].Status == "NEW" || lst[1].Status == "PROCESSING")
	assert.Equal(gs.T(), orderNum2, lst[1].Number)
	assert.True(gs.T(), lst[1].Status == "NEW" || lst[1].Status == "PROCESSING")

	// Register another user
	userLogin, userPassword = uuid.NewString(), uuid.NewString()

	resp, err = gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/register")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())

	// Get empty orders list
	resp, err = gs.httpClient.R().
		Get("/api/user/orders")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusNoContent, resp.StatusCode())

	// Add order, which exists from another user
	resp, err = gs.httpClient.R().
		SetBody(orderNum1).
		SetHeader("Content-Type", "text/plain").
		Post("/api/user/orders")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusConflict, resp.StatusCode())
}
