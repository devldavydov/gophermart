package gophermart

import "github.com/go-resty/resty/v2"

func (gs *GophermartSuite) userRegister(userLogin, userPassword string) (*resty.Response, error) {
	return gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/register")
}

func (gs *GophermartSuite) userLogin(userLogin, userPassword string) (*resty.Response, error) {
	return gs.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/login")
}

func (gs *GophermartSuite) userAddOrder(orderNum string) (*resty.Response, error) {
	return gs.httpClient.R().
		SetBody(orderNum).
		SetHeader("Content-Type", "text/plain").
		Post("/api/user/orders")
}
