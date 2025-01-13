/*
Пакет logger содержит объявление интерфейса логгера
*/
package logger

type Logger interface {
	// Debug логирование на уровне DEBUG
	Debug(msg string, args ...interface{})
	// Info логирование на уровне INFO
	Info(msg string, args ...interface{})
	// Warn логирование на уровне WARN
	Warn(msg string, args ...interface{})
	// Error логирование на уровне ERROR
	Error(msg string, args ...interface{})
}
