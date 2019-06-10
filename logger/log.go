package logger

import "gopkg.in/jucardi/go-logger-lib.v1/log"

const LoggerName = "jucardi/go-db"

var logger ILogger = log.Get(LoggerName)

func Get() ILogger {
	return logger
}

func Set(l ILogger) {
	logger = l
}

type ILogger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}
