package gophermart

import (
	"net/http"

	"github.com/google/uuid"
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
		gs.NoError(err)
		gs.Equal(http.StatusUnauthorized, resp.StatusCode())
	}
}

func (gs *GophermartSuite) TestRegisterLoginLogout() {
	gs.Run("register with wrong request", func() {
		resp, err := gs.httpClient.R().
			SetBody("foobar").
			Post("/api/user/register")
		gs.NoError(err)
		gs.Equal(http.StatusBadRequest, resp.StatusCode())
	})

	gs.Run("register successful", func() {
		userLogin, userPassword := uuid.NewString(), uuid.NewString()

		resp := gs.userRegister(userLogin, userPassword, http.StatusOK)

		authCookie := resp.Header().Get("Set-Cookie")
		gs.NotEqual("", authCookie)
	})

	gs.Run("register same user twice", func() {
		userLogin, userPassword := uuid.NewString(), uuid.NewString()

		gs.userRegister(userLogin, userPassword, http.StatusOK)
		gs.userRegister(userLogin, userPassword, http.StatusConflict)
	})

	gs.Run("register and check url", func() {
		userLogin, userPassword := uuid.NewString(), uuid.NewString()

		gs.userRegister(userLogin, userPassword, http.StatusOK)
		gs.userGetBalance(http.StatusOK)
	})

	gs.Run("register, logout and check url", func() {
		userLogin, userPassword := uuid.NewString(), uuid.NewString()

		gs.userRegister(userLogin, userPassword, http.StatusOK)
		gs.userLogout()

		gs.userGetBalance(http.StatusUnauthorized)
	})

	gs.Run("login with wrong request", func() {
		resp, err := gs.httpClient.R().
			SetBody("foobar").
			Post("/api/user/login")
		gs.NoError(err)
		gs.Equal(http.StatusBadRequest, resp.StatusCode())
	})

	gs.Run("register and login with wrong creds", func() {
		userLogin, userPassword := uuid.NewString(), uuid.NewString()

		gs.userRegister(userLogin, userPassword, http.StatusOK)
		gs.userLogin("foo", "bar", http.StatusUnauthorized)
	})

	gs.Run("register, logout, correct login and check url", func() {
		userLogin, userPassword := uuid.NewString(), uuid.NewString()

		gs.userRegister(userLogin, userPassword, http.StatusOK)
		gs.userLogout()
		gs.userLogin(userLogin, userPassword, http.StatusOK)

		gs.userGetBalance(http.StatusOK)
	})
}
