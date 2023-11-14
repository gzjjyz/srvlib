package redis

import (
	"context"

	"github.com/995933447/log-go"
	"github.com/gzjjyz/logger"
)

var Logger *log.Logger

var _ log.LoggerWriter = (*LoggerWriter)(nil)

type LoggerWriter struct {
}

func (l LoggerWriter) Write(ctx context.Context, level log.Level, format string, args ...interface{}) error {
	switch level {
	case log.LevelDebug:
		logger.LogDebug(format, args)
	case log.LevelInfo:
		logger.LogInfo(format, args)
	case log.LevelWarn:
		logger.LogWarn(format, args)
	case log.LevelError:
		logger.LogError(format, args)
	case log.LevelFatal:
		logger.LogFatal(format, args)
	case log.LevelPanic:
		logger.LogStack(format, args)
	}
	return nil
}

func (l LoggerWriter) Flush() error {
	logger.Flush()
	return nil
}

func init() {
	Logger = log.NewLogger(&LoggerWriter{})
}
