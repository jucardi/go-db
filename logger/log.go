package logger

import "gopkg.in/jucardi/go-logger-lib.v1/log"

const LoggerName = "jucardi/go-db"

func Get() log.ILogger {
	return log.Get(LoggerName)
}

func Set(logger log.ILogger) log.ILogger {
	return log.Register(LoggerName, logger)
}
