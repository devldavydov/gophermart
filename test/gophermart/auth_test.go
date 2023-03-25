package gophermart

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func (gs *GophermartSuite) TestApiWithoutAuth() {
	for _, api := range []struct {
		url    string
		method string
	}{
		{"/api/user/balance", http.MethodGet},
		{"/api/user/balance/withdraw", http.MethodPost},
		{"/api/user/withdrawals", http.MethodGet},
		{"/api/user/logout", http.MethodPost},
		{"/api/user/orders", http.MethodPost},
		{"/api/user/orders", http.MethodGet},
	} {
		resp, err := gs.httpClient.R().Execute(api.method, api.url)
		assert.NoError(gs.T(), err)
		assert.Equal(gs.T(), http.StatusUnauthorized, resp.StatusCode())
	}
}

func (gs *GophermartSuite) TestRegisterLoginLogout() {
	userLogin, userPassword := uuid.NewString(), uuid.NewString()

	// Register with wrong request
	resp, err := gs.httpClient.R().
		SetBody("foobar").
		Post("/api/user/register")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusBadRequest, resp.StatusCode())

	// Register
	resp, err = gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/register")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())

	authCookie := resp.Header().Get("Set-Cookie")
	assert.NotEqual(gs.T(), "", authCookie)

	// Register same user twice
	resp, err = gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/register")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusConflict, resp.StatusCode())

	// Try url after registered
	resp, err = gs.httpClient.R().Get("/api/user/balance")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())

	// Logout
	resp, err = gs.httpClient.R().Post("/api/user/logout")
	assert.NoError(gs.T(), err)

	// Try url after logout
	resp, err = gs.httpClient.R().Get("/api/user/balance")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusUnauthorized, resp.StatusCode())

	// Login with wrong request
	resp, err = gs.httpClient.R().
		SetBody("foobar").
		Post("/api/user/login")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusBadRequest, resp.StatusCode())

	// Login with wrong credentials
	resp, err = gs.httpClient.R().
		SetBody(userAuth{Login: "foo", Password: "bar"}).
		Post("/api/user/login")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusUnauthorized, resp.StatusCode())

	// Login
	resp, err = gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/login")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())

	// Try url after login
	resp, err = gs.httpClient.R().Get("/api/user/balance")
	assert.NoError(gs.T(), err)
	assert.Equal(gs.T(), http.StatusOK, resp.StatusCode())
}
