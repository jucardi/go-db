package mgo

import "gopkg.in/jucardi/go-logger-lib.v1/log"

type mgoLogger struct {
	logger log.ILogger
}

func (l *mgoLogger) Output(calldepth int, s string) error {
	l.logger.Debug(s)
	return nil
}

func wrapLogger(logger log.ILogger) *mgoLogger {
	return &mgoLogger{logger: logger}
}