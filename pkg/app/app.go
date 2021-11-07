package app

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
