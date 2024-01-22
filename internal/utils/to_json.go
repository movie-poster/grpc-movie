package utils

import (
	"encoding/json"

	"go.uber.org/zap"
)

func DataToJSON(data ...interface{}) string {
	marshal, err := json.Marshal(data)
	if err != nil {
		LoggerZap.Error("No se puede codificar la data", zap.Any("Data:", data))
	}
	return string(marshal)
}

func DataMapToJSON(data map[string]interface{}) string {
	marshal, err := json.Marshal(data)
	if err != nil {
		LoggerZap.Error("No se puede codificar la data", zap.Any("Data:", data))
	}
	return string(marshal)
}
