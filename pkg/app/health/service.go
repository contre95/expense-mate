package health

import (
	"expenses-app/pkg/app"
)

type Service struct {
	logger []app.Logger // Whynot?
}

// NewService returns a new HealthHandler
func NewService(l ...app.Logger) Service {
	if l != nil {
		return Service{logger: l}
	}
	return Service{}
}

func (s *Service) Ping() string {
	for _, log := range s.logger {
		log.Info("Healcheck ok")
	}
	return "pong"
}
