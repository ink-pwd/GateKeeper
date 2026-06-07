// logger/logger.go
package logger

import (
	"log"
	"os"
)

// Logger - интерфейс, который можно реализовать разными логгерами
type Logger interface {
	Info(format string, args ...any)
	Error(format string, args ...any)
	Fatal(format string, args ...any)
}

// StdLogger - реализация на стандартном log.Logger
type StdLogger struct {
	logger *log.Logger
}

// NewStdLogger создает новый StdLogger
func NewStdLogger() *StdLogger {
	return &StdLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *StdLogger) Info(format string, args ...any) {
	l.logger.Printf("[INFO] "+format, args...)
}

func (l *StdLogger) Error(format string, args ...any) {
	l.logger.Printf("[ERROR] "+format, args...)
}

func (l *StdLogger) Fatal(format string, args ...any) {
	l.logger.Fatalf("[FATAL] "+format, args...)
}
