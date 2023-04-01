package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/devldavydov/gophermart/internal/gophermart"
)

const (
	_defaultLogLevel             = "DEBUG"
	_defaultRunAddress           = "127.0.0.1:8080"
	_defaultDatabaseDsn          = ""
	_defaultAccrualSystemAddress = "http://127.0.0.1:9090"
	_defaultSessionSecret        = "secret"
	_defaultShutdownTimeout      = 10 * time.Second
	_defaultAccrualThreadNum     = 2
	_defaultOrderDBScanInterval  = 1 * time.Second
)

type Config struct {
	LogLevel             string        `env:"LOG_LEVEL"`
	RunAddress           string        `env:"RUN_ADDRESS"`
	DatabaseDsn          string        `env:"DATABASE_URI"`
	AccrualSystemAddress string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SessionSecret        string        `env:"SESSION_SECRET"`
	ShutdownTimeout      time.Duration `env:"SHUTDOWN_TIMEOUT"`
	AccrualThreadNum     int           `env:"ACCRUAL_THREAD_NUM"`
	OrderDBScanInterval  time.Duration `env:"ORDER_DB_SCAN_INTERVAL"`
}

func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	var err error
	config := &Config{}

	// Check flags
	flagSet.StringVar(&config.LogLevel, "l", _defaultLogLevel, "log level")
	flagSet.StringVar(&config.RunAddress, "a", _defaultRunAddress, "run address")
	flagSet.StringVar(&config.DatabaseDsn, "d", _defaultDatabaseDsn, "database uri")
	flagSet.StringVar(&config.AccrualSystemAddress, "r", _defaultAccrualSystemAddress, "accrual system address")
	flagSet.StringVar(&config.SessionSecret, "s", _defaultSessionSecret, "session secret")
	flagSet.DurationVar(&config.ShutdownTimeout, "t", _defaultShutdownTimeout, "shutdown timeout")
	flagSet.IntVar(&config.AccrualThreadNum, "n", _defaultAccrualThreadNum, "accrual thread num")
	flagSet.DurationVar(&config.OrderDBScanInterval, "o", _defaultOrderDBScanInterval, "order db scan interval")

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flagSet.PrintDefaults()
	}
	err = flagSet.Parse(flags)
	if err != nil {
		return nil, err
	}

	// Check env
	if err = env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}

func ServiceSettingsAdapt(config *Config) (*gophermart.ServiceSettings, error) {
	serviceSettings, err := gophermart.NewServiceSettings(
		"http://"+config.RunAddress,
		config.DatabaseDsn,
		config.AccrualSystemAddress,
		config.SessionSecret,
		config.ShutdownTimeout,
		config.AccrualThreadNum,
		config.OrderDBScanInterval,
	)
	if err != nil {
		return nil, err
	}
	return serviceSettings, nil
}
