package logger

import (
	"log/slog"
	"os"
)

type SlogLogger struct {
	slog *slog.Logger
}

func New(level string) *SlogLogger {
	var slogLevel slog.Level
	switch level {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "INFO":
		slogLevel = slog.LevelInfo
	case "WARN":
		slogLevel = slog.LevelWarn
	case "ERROR":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelWarn
	}
	opts := &slog.HandlerOptions{Level: slogLevel}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	return &SlogLogger{slog: logger}
}

func (l *SlogLogger) Debug(msg string, args ...interface{}) {
	l.slog.Debug(msg, args...)
}

func (l *SlogLogger) Info(msg string, args ...interface{}) {
	l.slog.Info(msg, args...)
}

func (l *SlogLogger) Warn(msg string, args ...interface{}) {
	l.slog.Warn(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...interface{}) {
	l.slog.Error(msg, args...)
}
