package helper

import (
	"log"
	"os"
	"project_pos_app/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog(cfg config.Config) (*zap.Logger, error) {

	logLevel := zap.InfoLevel
	if cfg.Debug {
		logLevel = zap.DebugLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		MessageKey:     "M",
		CallerKey:      "C",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	file, err := os.Create("app.log")
	if err != nil {
		log.Panicf("Failed to open log file: %v", err)
		return nil, err
	}

	writer := zapcore.AddSync(file)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		logLevel,
	)
	logger := zap.New(core)
	logger.Info("Logger initialized successfully")

	return logger, nil
}
