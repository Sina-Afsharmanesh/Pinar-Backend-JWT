package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logFilePath    = "default.log"
	critFilePath   = "critical.log"
	serverFilePath = "server.log"
)

var (
	Logger     = createLogger(logFilePath, false)
	ErrLogger  = createLogger(critFilePath, true)
	ServLogger = createLogger(serverFilePath, false)
)

func createLogger(filePath string, isError bool) *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		LocalTime:  true,
		Filename:   filePath,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     90,
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	if isError {
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	productionCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(productionCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core)
}
