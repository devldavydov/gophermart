package gophermart

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	httpClient *resty.Client
}

func NewClient(gophermartSrvAddress string) *Client {
	jar, _ := cookiejar.New(nil)
	httpClient := resty.New().
		SetBaseURL(gophermartSrvAddress).
		SetCookieJar(jar)

	return &Client{httpClient: httpClient}
}

func (c *Client) UserRegister(userLogin, userPassword string) error {
	resp, err := c.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/register")
	if err != nil {
		return err
	}

	if err = checkCommonError(resp.StatusCode()); err != nil {
		return err
	}

	if resp.StatusCode() == http.StatusConflict {
		return ErrUserAlreadyExists
	}

	return nil
}

func (c *Client) UserLogin(userLogin, userPassword string) error {
	resp, err := c.httpClient.R().
		SetBody(userAuth{Login: userLogin, Password: userPassword}).
		Post("/api/user/login")
	if err != nil {
		return err
	}

	if err = checkCommonError(resp.StatusCode()); err != nil {
		return err
	}

	return nil
}

func (c *Client) UserLogout() error {
	resp, err := c.httpClient.R().Post("/api/user/logout")
	if err != nil {
		return err
	}

	if err = checkCommonError(resp.StatusCode()); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetBalance() (*UserBalance, error) {
	var blnc UserBalance
	resp, err := c.httpClient.R().
		SetResult(&blnc).
		Get("/api/user/balance")
	if err != nil {
		return nil, err
	}

	if err = checkCommonError(resp.StatusCode()); err != nil {
		return nil, err
	}

	return &blnc, nil
}

func (c *Client) GetBalanceWithdrawals() (UserBalanceWithdrawals, error) {
	var wthdrwls UserBalanceWithdrawals
	resp, err := c.httpClient.R().
		SetResult(&wthdrwls).
		Get("/api/user/withdrawals")
	if err != nil {
		return nil, err
	}

	if err = checkCommonError(resp.StatusCode()); err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil, ErrNoBalanceWithdrawals
	}

	return wthdrwls, nil
}

func (c *Client) BalanceWithdraw(orderNum string, sum float64) error {
	resp, err := c.httpClient.R().
		SetBody(userBalanceWithdraw{Order: orderNum, Sum: sum}).
		Post("/api/user/balance/withdraw")
	if err != nil {
		return err
	}

	if err = checkCommonError(resp.StatusCode()); err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusUnprocessableEntity:
		return ErrBalanceWithdrawWrongFormat
	case http.StatusPaymentRequired:
		return ErrBalanceWithdrawPaymentRequired
	default:
		return nil
	}
}

func (c *Client) AddOrder(orderNum string) error {
	resp, err := c.httpClient.R().
		SetBody(orderNum).
		SetHeader("Content-Type", "text/plain").
		Post("/api/user/orders")
	if err != nil {
		return err
	}

	if err = checkCommonError(resp.StatusCode()); err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusUnprocessableEntity:
		return ErrOrderWrongFormat
	case http.StatusConflict:
		return ErrOrderAlreadyExists
	case http.StatusOK:
		return ErrOrderAlreadyAccepted
	default:
		return nil
	}
}

func (c *Client) GetOrders() (OrderList, error) {
	var lst OrderList
	resp, err := c.httpClient.R().
		SetResult(&lst).
		Get("/api/user/orders")

	if err != nil {
		return nil, err
	}

	if err = checkCommonError(resp.StatusCode()); err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil, ErrNoOrders
	}

	return lst, nil
}

func checkCommonError(statusCode int) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusInternalServerError:
		return ErrInternalError
	default:
		return nil
	}
}
