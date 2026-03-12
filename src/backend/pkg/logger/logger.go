package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log 全局日志实例
	Log *zap.Logger
)

// InitLogger 初始化日志
func InitLogger(mode string) error {
	var config zap.Config

	if mode == "release" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	log, err := config.Build()
	if err != nil {
		return err
	}

	Log = log
	return nil
}

// GetLogger 获取全局日志实例
func GetLogger() *zap.Logger {
	return Log
}
