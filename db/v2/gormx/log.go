package gormx

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/gzjjyz/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Logger struct {
	slowThreshold time.Duration
}

func NewLogger(slowThreshold time.Duration) *Logger {
	return &Logger{
		slowThreshold: slowThreshold,
	}
}

const (
	traceStr     = "%s [%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s [%.3fms] [rows:%v] %s"
)

func (l *Logger) LogMode(_ logger.LogLevel) logger.Interface {
	return l
}

func (l *Logger) Info(_ context.Context, format string, args ...interface{}) {
	log.LogInfo(format, args...)
}

func (l *Logger) Warn(_ context.Context, format string, args ...interface{}) {
	log.LogWarn(format, args...)
}

func (l *Logger) Error(_ context.Context, format string, args ...interface{}) {
	log.LogError(format, args...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			l.Error(ctx, traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Error(ctx, traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.slowThreshold && l.slowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
		if rows == -1 {
			l.Warn(ctx, traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Warn(ctx, traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			l.Info(ctx, traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Info(ctx, traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
