package utils

import (
	"github.com/gzjjyz/logger"
)

func SafeLogErr(err error, printWhileLoggerNoReady bool) {
	logger.LogError(err.Error())
}

func SafeLogWarn(printWhileLoggerNoReady bool, format string, args ...interface{}) {
	logger.LogWarn(format, args...)
}
