// Package zap provides a Logger that writes to zap.SugaredLogger.
package zap

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapSugarLogger struct {
	logger     *zap.SugaredLogger
	defaultLog func(kv ...interface{}) error
}

func (l zapSugarLogger) Log(kv ...interface{}) error {

	if len(kv) <= 1 {
		l.logger.Info(kv...)
		return nil
	}

	if kv[0] == level.Key() {
		switch kv[1] {
		case level.ErrorValue():
			l.logger.Errorw("", kv[2:]...)
		case level.WarnValue():
			l.logger.Warnw("", kv[2:]...)
		case level.InfoValue():
			l.logger.Infow("", kv[2:]...)
		case level.DebugValue():
			l.logger.Debugw("", kv[2:]...)
		default:
			_ = l.defaultLog(kv...)
		}
	} else {
		_ = l.defaultLog(kv...)
	}
	return nil
}

// NewZapSugarLogger returns a Go kit log.Logger that sends
// log events to a zap.Logger. if no level provided with the
// log event, the param level is used.
func NewZapSugarLogger(logger *zap.Logger, level zapcore.Level) log.Logger {
	sugarLogger := logger.WithOptions(zap.AddCallerSkip(2)).Sugar()

	var defaultLog func(msg string, keysAndValues ...interface{})

	switch level {
	case zapcore.DebugLevel:
		defaultLog = sugarLogger.Debugw
	case zapcore.InfoLevel:
		defaultLog = sugarLogger.Infow
	case zapcore.WarnLevel:
		defaultLog = sugarLogger.Warnw
	case zapcore.ErrorLevel:
		defaultLog = sugarLogger.Errorw
	case zapcore.DPanicLevel:
		defaultLog = sugarLogger.DPanicw
	case zapcore.PanicLevel:
		defaultLog = sugarLogger.Panicw
	case zapcore.FatalLevel:
		defaultLog = sugarLogger.Fatalw
	default:
		defaultLog = sugarLogger.Infow
	}

	return zapSugarLogger{
		logger: sugarLogger,
		defaultLog: func(kv ...interface{}) error {
			defaultLog("", kv...)
			return nil
		},
	}
}
