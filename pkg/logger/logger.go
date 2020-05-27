package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var VerboseLogging bool = false
var ZapLogger *zap.SugaredLogger = getLogger()

func Debug(msg string) {
    ZapLogger.Debug(msg)
}

func Info(msg string) {
    ZapLogger.Info(msg)
}

func Warn(msg string) {
    ZapLogger.Warn(msg)
}

func Error(msg string) {
    ZapLogger.Error(msg)
}

func Fatal(msg string) {
    ZapLogger.Fatal(msg)
}

func SetVerbose(verbose bool) {
    VerboseLogging = verbose
    ZapLogger = getLogger()
}

func getLogger() *zap.SugaredLogger {
    config := zap.NewDevelopmentConfig()
    config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    if !VerboseLogging {
        config = zap.NewProductionConfig()
    }
    logger, _ := config.Build()
    return logger.Sugar()
}
