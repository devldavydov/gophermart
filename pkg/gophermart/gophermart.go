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
