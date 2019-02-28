package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

type Sugar struct {
	*zap.SugaredLogger
}

type LoggerOpt func(config *zap.Config)

func NewLogger(opts ...LoggerOpt) *Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = TimeEncoder

	for _, opt := range opts {
		opt(&config)
	}
	// config.Level = NewAtomicLevelAt(DebugLevel)

	logger, _ := config.Build()
	return &Logger{logger}
}

func (l *Logger) SugarLogger() *Sugar {
	return NewSugar(l.Sugar())
}

func NewSugar(sugar *zap.SugaredLogger) *Sugar {
	return &Sugar{sugar}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

var (
	logger *Logger
	sugar  *Sugar
)

func L() *Logger {
	return logger
}

func S() *Sugar {
	return sugar
}
func RegistryLogger(l *Logger) {
	logger = l
}

func RegistrySugar(s *Sugar) {
	sugar = s
}

func init() {
	logger = NewLogger(Level(zapcore.ErrorLevel))
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.SugarLogger()
}
