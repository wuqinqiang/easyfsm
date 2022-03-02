package log

import (
	"log"
)

// DefaultLogger is default logger.
var DefaultLogger Logger = NewStdLogger(log.Writer())

// Logger is a logger interface.
type Logger interface {
	Log(level Level, keyvals ...interface{})
}

func SetLogger(logger Logger) {
	DefaultLogger = logger
}
