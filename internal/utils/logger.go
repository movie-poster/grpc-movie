package utils

import (
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LoggerZap *zap.Logger

func SetupLoggerZap() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("02/Jan/2006:15:04:05 -0700")
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, _ := os.OpenFile("../internal/logs/log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	LoggerZap = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func LogError(message, description, method string, status int, data ...interface{}) {
	m := DataToJSON(data)
	LoggerZap.Error(
		message,
		zap.String("Description:", description),
		zap.String("Method:", method),
		zap.Int("Status:", status),
		zap.String("Data:", m),
	)
}

func LogWarning(message, description, method string, data ...interface{}) {
	m := DataToJSON(data)
	LoggerZap.Warn(
		message,
		zap.String("Description:", description),
		zap.String("Method:", method),
		zap.String("Data:", m),
	)
}

func LogSuccess(message, method string, data ...interface{}) {
	m := DataToJSON(data)
	LoggerZap.Info(
		message,
		zap.String("Method:", method),
		zap.Int("Status:", http.StatusOK),
		zap.String("Data:", m),
	)
}
