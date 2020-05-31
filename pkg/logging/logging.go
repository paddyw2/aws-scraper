package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Logger *zap.SugaredLogger
	Level  int
}

func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(args)
}

func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(args)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Warn(args)
}

func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(args)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args)
}

func NewLogger(level int) *Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	var loggingLevel zapcore.Level
	if level == 0 {
		loggingLevel = zap.ErrorLevel
	} else if level == 1 {
		loggingLevel = zap.InfoLevel
	} else {
		loggingLevel = zap.DebugLevel
	}
	config.Level.SetLevel(loggingLevel)
	logger, _ := config.Build()
	newLogger := Logger{Logger: logger.Sugar(), Level: level}
	return &newLogger
}
