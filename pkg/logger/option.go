package logger

import (
	zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Level(lvl zapcore.Level) LoggerOpt {
	return func(config *zap.Config) {
		config.Level = zap.NewAtomicLevelAt(lvl)
	}
}
