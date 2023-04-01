package utils

import (
	"urlshortner/utils/logger"

	"go.uber.org/zap"
)

func HandlePanicGo() {
	if r := recover(); r != nil {
		defer logger.Log.Sync()
		logger.Log.Error("panic: ", zap.Any("error: ", r))
	}
}
