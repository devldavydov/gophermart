package gophermart

import (
	"net/url"
	"time"
)

type ServiceSettings struct {
	RunAddress           *url.URL
	DatabaseDsn          string
	AccrualSystemAddress *url.URL
	SessionSecret        string
	ShutdownTimeout      time.Duration
	AccrualThreadNum     int
	OrderDBScanInterval  time.Duration
}

func NewServiceSettings(
	runAddress, databaseDsn, accrualSystemAddress string,
	sessionSecret string, shutdownTimeout time.Duration,
	accrualThreadNum int, orderDBScanInterval time.Duration,
) (*ServiceSettings, error) {
	urlRunAddress, err := url.ParseRequestURI(runAddress)
	if err != nil {
		return nil, err
	}

	urlAccrualSystemAddress, err := url.ParseRequestURI(accrualSystemAddress)
	if err != nil {
		return nil, err
	}

	return &ServiceSettings{
		RunAddress:           urlRunAddress,
		DatabaseDsn:          databaseDsn,
		AccrualSystemAddress: urlAccrualSystemAddress,
		SessionSecret:        sessionSecret,
		ShutdownTimeout:      shutdownTimeout,
		AccrualThreadNum:     accrualThreadNum,
		OrderDBScanInterval:  orderDBScanInterval,
	}, nil
}
