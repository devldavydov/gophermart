package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/devldavydov/gophermart/internal/common/env"
	"github.com/devldavydov/gophermart/internal/gophermart"
)

const (
	_defaultLogLevel             = "DEBUG"
	_defaultRunAddress           = "127.0.0.1:8080"
	_defaultDatabaseDsn          = ""
	_defaultAccrualSystemAddress = "http://127.0.0.1:9090"
	_defaultSessionSecret        = "secret"
	_defaultShutdownTimeout      = 10 * time.Second
)

type Config struct {
	LogLevel             string
	RunAddress           string
	DatabaseDsn          string
	AccrualSystemAddress string
	SessionSecret        string
	ShutdownTimeout      time.Duration
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

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flagSet.PrintDefaults()
	}
	err = flagSet.Parse(flags)
	if err != nil {
		return nil, err
	}

	// Check env
	config.LogLevel, err = env.GetVariable("LOG_LEVEL", env.CastString, config.LogLevel)
	if err != nil {
		return nil, err
	}

	config.RunAddress, err = env.GetVariable("RUN_ADDRESS", env.CastString, config.RunAddress)
	if err != nil {
		return nil, err
	}

	config.DatabaseDsn, err = env.GetVariable("DATABASE_URI", env.CastString, config.DatabaseDsn)
	if err != nil {
		return nil, err
	}

	config.AccrualSystemAddress, err = env.GetVariable("ACCRUAL_SYSTEM_ADDRESS", env.CastString, config.AccrualSystemAddress)
	if err != nil {
		return nil, err
	}

	config.SessionSecret, err = env.GetVariable("SESSION_SECRET", env.CastString, config.SessionSecret)
	if err != nil {
		return nil, err
	}

	config.ShutdownTimeout, err = env.GetVariable("SHUTDOWN_TIMEOUT", env.CastDuration, config.ShutdownTimeout)
	if err != nil {
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
	)
	if err != nil {
		return nil, err
	}
	return serviceSettings, nil
}
