package gophermart

import (
	"context"
	"fmt"
	"net/http"

	"github.com/devldavydov/gophermart/internal/gophermart/handler"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Service struct {
	settings *ServiceSettings
	logger   *logrus.Logger
}

func NewService(settings *ServiceSettings, logger *logrus.Logger) *Service {
	return &Service{settings: settings, logger: logger}
}

func (s *Service) Start(ctx context.Context) error {
	s.logger.Infof("Service started on [%s]", s.settings.RunAddress)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	handler.Init(router)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.settings.RunAddress.Hostname(), s.settings.RunAddress.Port()),
		Handler: router,
	}

	errChan := make(chan error)
	go func(ch chan error) {
		ch <- httpServer.ListenAndServe()
	}(errChan)

	select {
	case err := <-errChan:
		return fmt.Errorf("service exited with err: %w", err)
	case <-ctx.Done():
		s.logger.Infof("Service context canceled")

		ctx, cancel := context.WithTimeout(context.Background(), s.settings.ShutdownTimeout)
		defer cancel()

		err := httpServer.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("service shutdown err: %w", err)
		}

		s.logger.Info("Service finished")
		return nil
	}
}
