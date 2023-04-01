package gophermart

import (
	"net/http"

	"github.com/devldavydov/gophermart/pkg/gophermart"
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
		gs.NoError(gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))
	})

	gs.Run("register same user twice", func() {
		userLogin, userPassword := uuid.NewString(), uuid.NewString()

		gs.NoError(gs.gCli.UserRegister(userLogin, userPassword))
		gs.ErrorIs(gophermart.ErrUserAlreadyExists, gs.gCli.UserRegister(userLogin, userPassword))
	})

	gs.Run("register and check url", func() {
		gs.NoError(gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))

		_, err := gs.gCli.GetBalance()
		gs.NoError(err)
	})

	gs.Run("register, logout and check url", func() {
		gs.NoError(gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))
		gs.NoError(gs.gCli.UserLogout())

		_, err := gs.gCli.GetBalance()
		gs.ErrorIs(gophermart.ErrUnauthorized, err)
	})

	gs.Run("login with wrong request", func() {
		resp, err := gs.httpClient.R().
			SetBody("foobar").
			Post("/api/user/login")
		gs.NoError(err)
		gs.Equal(http.StatusBadRequest, resp.StatusCode())
	})

	gs.Run("register and login with wrong creds", func() {
		gs.NoError(gs.gCli.UserRegister(uuid.NewString(), uuid.NewString()))
		gs.ErrorIs(gophermart.ErrUnauthorized, gs.gCli.UserLogin("foo", "bar"))
	})

	gs.Run("register, logout, correct login and check url", func() {
		userLogin, userPassword := uuid.NewString(), uuid.NewString()

		gs.NoError(gs.gCli.UserRegister(userLogin, userPassword))
		gs.NoError(gs.gCli.UserLogout())
		gs.NoError(gs.gCli.UserLogin(userLogin, userPassword))

		_, err := gs.gCli.GetBalance()
		gs.NoError(err)
	})
}
