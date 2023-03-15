package gophermart

import (
	"context"

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
	return nil
}
