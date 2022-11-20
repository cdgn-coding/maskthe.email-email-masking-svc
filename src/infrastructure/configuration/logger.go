package configuration

import (
	"go.uber.org/zap"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type logger struct {
	sugarLogger *zap.SugaredLogger
}

func NewLogger(level string) *logger {
	switch level {
	case "production":
		zapLogger, _ := zap.NewProduction()
		sugarLogger := zapLogger.Sugar()
		return &logger{sugarLogger}
	default:
		zapLogger, _ := zap.NewDevelopment()
		sugarLogger := zapLogger.Sugar()
		return &logger{sugarLogger}
	}
}

func (l logger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l logger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l logger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l logger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l logger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}
