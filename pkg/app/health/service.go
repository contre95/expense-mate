package health

import (
	"expenses-app/pkg/app"
)

type Service struct {
	logger app.Logger
}

// NewService returns a new HealthHandler
func NewService(l app.Logger) Service {
	return Service{l}
}

func (s *Service) Ping() string {
	s.logger.Debug("Ping requested")
	return "pong"
}
