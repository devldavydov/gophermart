package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_log "github.com/devldavydov/gophermart/internal/common/log"
	"github.com/devldavydov/gophermart/internal/gophermart"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := LoadConfig(*flag.CommandLine, os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to load flag and ENV settings: %w", err)
	}

	logger, err := _log.NewLogger(config.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	serviceSettings, err := ServiceSettingsAdapt(config)
	if err != nil {
		return fmt.Errorf("failed to create service settings: %w", err)
	}

	serviceSettings.DatabaseDsn = "postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable"
	service := gophermart.NewService(serviceSettings, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return service.Start(ctx)
}
