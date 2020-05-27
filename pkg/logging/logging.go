package logging

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger(verbose bool) *zap.SugaredLogger {
    config := zap.NewDevelopmentConfig()
    config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    if !verbose {
        config = zap.NewProductionConfig()
    }
    logger, _ := config.Build()
    return logger.Sugar()
}
