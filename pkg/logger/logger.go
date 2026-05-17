package logger

import "context"

type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)

	DBDebug(ctx context.Context, msg string, keysAndValues ...any)
	DBInfo(ctx context.Context, msg string, keysAndValues ...any)
	DBWarn(ctx context.Context, msg string, keysAndValues ...any)
	DBError(ctx context.Context, msg string, keysAndValues ...any)
}

// ConsoleLogger is a simple logger that writes to stdout
type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Debug(msg string, keysAndValues ...any)                        {}
func (l *ConsoleLogger) Info(msg string, keysAndValues ...any)                         {}
func (l *ConsoleLogger) Warn(msg string, keysAndValues ...any)                         {}
func (l *ConsoleLogger) Error(msg string, keysAndValues ...any)                        {}
func (l *ConsoleLogger) DBDebug(ctx context.Context, msg string, keysAndValues ...any) {}
func (l *ConsoleLogger) DBInfo(ctx context.Context, msg string, keysAndValues ...any)  {}
func (l *ConsoleLogger) DBWarn(ctx context.Context, msg string, keysAndValues ...any)  {}
func (l *ConsoleLogger) DBError(ctx context.Context, msg string, keysAndValues ...any) {}
