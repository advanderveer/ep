package ep

import (
	"log"
)

// Logger interface may be implemented to allow the endpoints to provide
// feedback in what ever way is preferred
type Logger interface {

	// LogServerErrRender is called when the response will render a server error
	LogServerErrRender(err error)

	// LogClientErrRender is called when the response will render a client error
	LogClientErrRender(err error)

	// LogAppErrRender is called when the response will be a rendered app errror
	LogAppErrRender(err *AppError)
}

// StdLogger creates a logger using the standard library logging package
type StdLogger struct{ l *log.Logger }

// NewStdLogger inits a new standard library logger. If logs is nil it will
// use the global logs.* methods for printing
func NewStdLogger(logs *log.Logger) *StdLogger {
	return &StdLogger{logs}
}

func (l StdLogger) logf(f string, args ...interface{}) {
	if l.l != nil {
		l.l.Printf(f, args...)
	} else {
		log.Printf(f, args...)
	}
}

func (l StdLogger) LogServerErrRender(err error) {
	l.logf("ep: rendering server error: %v", err)
}
func (l StdLogger) LogClientErrRender(err error) {
	l.logf("ep: rendering client error: %v", err)
}

func (l StdLogger) LogAppErrRender(err *AppError) {
	l.logf("ep: rendering app error: %v", err)
}

// NopLogger is can be provided to disable logging
type NopLogger struct{}

func (l NopLogger) LogServerErrRender(err error)  {}
func (l NopLogger) LogClientErrRender(err error)  {}
func (l NopLogger) LogAppErrRender(err *AppError) {}
