package app

// "expenses-app/pkg/app/authenticating"
// "expenses-app/pkg/app/importing"
// "expenses-app/pkg/app/importing"
// "expenses-app/pkg/app/managing"

type Logger interface {
	// Use for logging informations
	Info(format string, i ...interface{})
	// Use for logging Warnings
	Warn(format string, i ...interface{})
	// Use for logging Errors
	Err(format string, i ...interface{})
	// Use for Debugging
	Debug(format string, i ...interface{})
}

type Hasher interface {
	// Hash hashes
	Hash(password string) (string, error)
	// CheckHash checks if a hash if string is equal a anotherone hashed
	CheckHash(password, hash string) bool
}
