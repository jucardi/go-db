package sql

import "gopkg.in/jucardi/go-logger-lib.v1/log"

type sqlLogger struct {
	logger log.ILogger
}

func (l *sqlLogger) Print(v ...interface{}) {
	l.logger.Info(v...)
}

func wrapLogger(logger log.ILogger) *sqlLogger {
	return &sqlLogger{logger: logger}
}