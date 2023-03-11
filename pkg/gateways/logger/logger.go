package logger

import (
	"fmt"
	"log"
)

type stdLogger struct {
	instanceName string
}

func NewSTDLogger(instanceName string, ctxColor string) *stdLogger {
	return &stdLogger{ctxColor + "[" + instanceName + "]" + RESET}
}

func (l *stdLogger) Info(format string, i ...interface{}) {
	fmt.Printf(l.instanceName)
	fmt.Printf(BLUE + BOLD + " [INFO] " + RESET)
	log.Printf(format, i...)
}

func (l *stdLogger) Warn(format string, i ...interface{}) {
	fmt.Printf(l.instanceName)
	fmt.Printf(YELLOW + BOLD + " [WARNING] " + RESET)
	log.Printf(format, i...)

}

func (l *stdLogger) Err(format string, i ...interface{}) {
	fmt.Printf(l.instanceName)
	fmt.Printf(RED + BOLD + " [ERROR] " + RESET)
	log.Printf(format, i...)

}

func (l *stdLogger) Debug(format string, i ...interface{}) {
	fmt.Printf(l.instanceName)
	fmt.Printf(PURPLE + BOLD + " [DEBUG] " + RESET)
	log.Printf(format, i...)
}

func (l *stdLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(l.instanceName)
	fmt.Printf(PURPLE + BOLD + " [-] " + RESET)
	log.Printf(format, v...)

}

func (l *stdLogger) Println(i ...interface{}) {
	fmt.Printf(l.instanceName)
	fmt.Printf(PURPLE + BOLD + " [-] " + RESET)
	log.Println(i...)
}
