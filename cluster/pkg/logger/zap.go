package logger

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewFileLogger creates a zap logger for file.
func NewFileLogger(level string, prefix int) *zap.Logger {
	var lvl zapcore.Level

	if err := lvl.Set(level); err != nil {
		log.Printf("cannot parse log level %s: %s", level, err)

		lvl = zapcore.WarnLevel
	}

	name := fmt.Sprintf("logs%d.csv", prefix)

	if err := os.Truncate(name, 0); err != nil {
		panic(err)
	}

	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	fileCore := zapcore.NewCore(encoder, zapcore.AddSync(file), lvl)
	cores := []zapcore.Core{
		fileCore,
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return logger
}
