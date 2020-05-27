package logging

import (
    "os"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

const LOG_LEVEL string = "GO_SCRAPER_LOG_LEVEL"

func GetLogger() *zap.SugaredLogger {
    logger := getLoggerByLevel()
    return logger.Sugar()
}

func getLoggerByLevel() *zap.Logger {
    debugLevel := os.Getenv(LOG_LEVEL)
    var config zap.Config
    config = zap.NewDevelopmentConfig()
    config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    if debugLevel == "PROD" {
        config = zap.NewProductionConfig()
    }
    logger, _ := config.Build()
    return logger
}
