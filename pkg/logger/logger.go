package logger

import (
    "os"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

const log_level string = "GO_SCRAPER_LOG_LEVEL"
var zapLogger *zap.SugaredLogger = getLogger()

func Debug(msg string) {
    zapLogger.Debug(msg)
}

func Info(msg string) {
    zapLogger.Info(msg)
}

func Warn(msg string) {
    zapLogger.Warn(msg)
}

func Error(msg string) {
    zapLogger.Error(msg)
}

func Fatal(msg string) {
    zapLogger.Fatal(msg)
}

func getLogger() *zap.SugaredLogger {
    logger := getLoggerByLevel()
    return logger.Sugar()
}

func getLoggerByLevel() *zap.Logger {
    debugLevel := os.Getenv(log_level)
    var config zap.Config
    config = zap.NewDevelopmentConfig()
    config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    if debugLevel == "PROD" {
        config = zap.NewProductionConfig()
    }
    logger, _ := config.Build()
    return logger
}
