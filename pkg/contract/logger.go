package contract

import (
	"context"
	gormlogger "gorm.io/gorm/logger"
)

// ILogger 定义了 onex 项目的日志接口. 该接口只包含了支持的日志记录方法.
type ILogger interface {
	Debugf(format string, args ...any)
	Debugw(msg string, keyvals ...any)

	Infof(format string, args ...any)
	Infow(msg string, keyvals ...any)

	Warnf(format string, args ...any)
	Warnw(msg string, keyvals ...any)

	Errorf(format string, args ...any)
	Errorw(err error, msg string, keyvals ...any)

	Panicf(format string, args ...any)
	Panicw(msg string, keyvals ...any)

	Fatalf(format string, args ...any)
	Fatalw(msg string, keyvals ...any)

	W(ctx context.Context) ILogger
	AddCallerSkip(skip int) ILogger

	Sync()

	gormlogger.Interface
}
