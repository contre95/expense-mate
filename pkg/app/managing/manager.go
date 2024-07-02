package managing

// Service just holds all the managing use cases
type Service struct {
	// UserCreator     UsersCreator
	CategoryDeleter   CategoryDeleter
	CategoryCreator   CategoryCreator
	CategoryUpdater   CategoryUpdater
	TelegramCommander TelegramCommander
	RuleManager       RuleManager
}

// NewService is the interctor for all Managing Use cases
func NewService(cd CategoryDeleter, cc CategoryCreator, cu CategoryUpdater, tc TelegramCommander, rm RuleManager) Service {
	return Service{cd, cc, cu, tc, rm}
}
