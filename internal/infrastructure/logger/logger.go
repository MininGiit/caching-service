/*
Пакет logger содержит имплементацию логера "cachingService/internal/logger"
*/
package logger

import (
	"log/slog"
	"os"
)

// SlogLogger структура логгра на основе slog
type SlogLogger struct {
	slog *slog.Logger
}

// New создание экземпляра логгера
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

// Debug логирование на уровне DEBUG
func (l *SlogLogger) Debug(msg string, args ...interface{}) {
	l.slog.Debug(msg, args...)
}

// Info логирование на уровне INFO
func (l *SlogLogger) Info(msg string, args ...interface{}) {
	l.slog.Info(msg, args...)
}

// Warn логирование на уровне WARN
func (l *SlogLogger) Warn(msg string, args ...interface{}) {
	l.slog.Warn(msg, args...)
}

// Error логирование на уровне ERROR
func (l *SlogLogger) Error(msg string, args ...interface{}) {
	l.slog.Error(msg, args...)
}
