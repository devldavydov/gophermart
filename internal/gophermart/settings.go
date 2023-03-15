package gophermart

import (
	"net/url"
	"time"
)

type ServiceSettings struct {
	RunAddress           *url.URL
	DatabaseDsn          string
	AccrualSystemAddress *url.URL
	ShutdownTimeout      time.Duration
}

func NewServiceSettings(runAddress, databaseDsn, accrualSystemAddress string, shutdownTimeout time.Duration) (*ServiceSettings, error) {
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
		ShutdownTimeout:      shutdownTimeout,
	}, nil
}
