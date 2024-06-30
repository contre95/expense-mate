package managing

// Service just holds all the managing use cases
type Service struct {
	// UserCreator     UsersCreator
	CategoryDeleter   CategoryDeleter
	CategoryCreator   CategoryCreator
	CategoryUpdater   CategoryUpdater
	TelegramCommander TelegramCommander
}

// NewService is the interctor for all Managing Use cases
func NewService(cd CategoryDeleter, cc CategoryCreator, cu CategoryUpdater, tc TelegramCommander) Service {
	return Service{cd, cc, cu, tc}
}
