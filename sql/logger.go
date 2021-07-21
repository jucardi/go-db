package sql

import "github.com/jucardi/go-logger-lib/log"

type sqlLogger struct {
	logger log.ILogger
}

func (l *sqlLogger) Print(v ...interface{}) {
	l.logger.Info(v...)
}

func wrapLogger(logger log.ILogger) *sqlLogger {
	return &sqlLogger{logger: logger}
}