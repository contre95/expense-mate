package health

import (
	"expenses-app/pkg/app"
	"sync/atomic"
)

type Service struct {
	logger    app.Logger
	botStatus *int32
}

// NewService returns a new HealthHandler
func NewService(l app.Logger, b *int32) Service {
	return Service{l, b}
}

func (s *Service) Ping() string {
	s.logger.Debug("Ping requested")
	return "pong"
}

func (s *Service) CheckBotHealth() string {
	s.logger.Debug("Bot status requested")
	if atomic.LoadInt32(s.botStatus) == 1 {
		return "running"
	}
	return "not running"
}
