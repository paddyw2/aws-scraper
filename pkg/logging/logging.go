package logging

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger(level int) *zap.SugaredLogger {
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
    return logger.Sugar()
}
