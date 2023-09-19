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
		logger.Debug(format, args)
	case log.LevelInfo:
		logger.Info(format, args)
	case log.LevelWarn:
		logger.Warn(format, args)
	case log.LevelError:
		logger.Errorf(format, args)
	case log.LevelFatal:
		logger.Fatalf(format, args)
	case log.LevelPanic:
		logger.Stack(format, args)
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
