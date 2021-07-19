package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)

	Core() zapcore.Core
	Named(s string) *zap.Logger
	Sync() error
	With(fields ...zap.Field) *zap.Logger
}
