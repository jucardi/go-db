package mgo

import "github.com/jucardi/go-logger-lib/log"

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