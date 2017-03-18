package gomolreplay

import (
	"os"

	"github.com/aphistic/gomol"
)

type (
	appExiter interface {
		Exit(code int)
	}

	osExiter struct{}
)

var (
	curExiter appExiter = &osExiter{}
)

func (exiter *osExiter) Exit(code int) {
	os.Exit(code)
}

func setExiter(exiter appExiter) {
	curExiter = exiter
}

// Dbg is a short-hand version of Debug
func (a *Adapter) Dbg(msg string) error {
	return a.Debug(msg)
}

// Dbgf is a short-hand version of Debugf
func (a *Adapter) Dbgf(msg string, args ...interface{}) error {
	return a.Debugf(msg, args...)
}

// Dbgm is a short-hand version of Debugm
func (a *Adapter) Dbgm(m *gomol.Attrs, msg string, args ...interface{}) error {
	return a.Debugm(m, msg, args...)
}

// Debug logs msg to all added loggers at LogLevel.LevelDebug
func (a *Adapter) Debug(msg string) error {
	return a.Log(gomol.LevelDebug, nil, msg)
}

/*
Debugf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelDebug
*/
func (a *Adapter) Debugf(msg string, args ...interface{}) error {
	return a.Log(gomol.LevelDebug, nil, msg, args...)
}

/*
Debugm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelDebug. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (a *Adapter) Debugm(m *gomol.Attrs, msg string, args ...interface{}) error {
	return a.Log(gomol.LevelDebug, m, msg, args...)
}

// Info logs msg to all added loggers at LogLevel.LevelInfo
func (a *Adapter) Info(msg string) error {
	return a.Log(gomol.LevelInfo, nil, msg)
}

/*
Infof uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelInfo
*/
func (a *Adapter) Infof(msg string, args ...interface{}) error {
	return a.Log(gomol.LevelInfo, nil, msg, args...)
}

/*
Infom uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelInfo. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (a *Adapter) Infom(m *gomol.Attrs, msg string, args ...interface{}) error {
	return a.Log(gomol.LevelInfo, m, msg, args...)
}

// Warn is a short-hand version of Warning
func (a *Adapter) Warn(msg string) error {
	return a.Warning(msg)
}

// Warnf is a short-hand version of Warningf
func (a *Adapter) Warnf(msg string, args ...interface{}) error {
	return a.Warningf(msg, args...)
}

// Warnm is a short-hand version of Warningm
func (a *Adapter) Warnm(m *gomol.Attrs, msg string, args ...interface{}) error {
	return a.Warningm(m, msg, args...)
}

/*
Warning uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning
*/
func (a *Adapter) Warning(msg string) error {
	return a.Log(gomol.LevelWarning, nil, msg)
}

/*
Warningf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning
*/
func (a *Adapter) Warningf(msg string, args ...interface{}) error {
	return a.Log(gomol.LevelWarning, nil, msg, args...)
}

/*
Warningm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (a *Adapter) Warningm(m *gomol.Attrs, msg string, args ...interface{}) error {
	return a.Log(gomol.LevelWarning, m, msg, args...)
}

// Err is a short-hand version of Error
func (a *Adapter) Err(msg string) error {
	return a.Error(msg)
}

// Errf is a short-hand version of Errorf
func (a *Adapter) Errf(msg string, args ...interface{}) error {
	return a.Errorf(msg, args...)
}

// Errm is a short-hand version of Errorm
func (a *Adapter) Errm(m *gomol.Attrs, msg string, args ...interface{}) error {
	return a.Errorm(m, msg, args...)
}

/*
Error uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError
*/
func (a *Adapter) Error(msg string) error {
	return a.Log(gomol.LevelError, nil, msg)
}

/*
Errorf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError
*/
func (a *Adapter) Errorf(msg string, args ...interface{}) error {
	return a.Log(gomol.LevelError, nil, msg, args...)
}

/*
Errorm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (a *Adapter) Errorm(m *gomol.Attrs, msg string, args ...interface{}) error {
	return a.Log(gomol.LevelError, m, msg, args...)
}

/*
Fatal uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal
*/
func (a *Adapter) Fatal(msg string) error {
	return a.Log(gomol.LevelFatal, nil, msg)
}

/*
Fatalf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal
*/
func (a *Adapter) Fatalf(msg string, args ...interface{}) error {
	return a.Log(gomol.LevelFatal, nil, msg, args...)
}

/*
Fatalm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (a *Adapter) Fatalm(m *gomol.Attrs, msg string, args ...interface{}) error {
	return a.Log(gomol.LevelFatal, m, msg, args...)
}

// Die will log a message using Fatal, call ShutdownLoggers and then exit the application with the provided exit code.
// This function is not subject to rollup and is always sent to the wrapped logger.
func (a *Adapter) Die(exitCode int, msg string) {
	a.Log(gomol.LevelFatal, nil, msg)
	a.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}

// Dief will log a message using Fatalf, call ShutdownLoggers and then exit the application with the provided exit code.
func (a *Adapter) Dief(exitCode int, msg string, args ...interface{}) {
	a.Log(gomol.LevelFatal, nil, msg, args...)
	a.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}

// Diem will log a message using Fatalm, call ShutdownLoggers and then exit the application with the provided exit code.
func (a *Adapter) Diem(exitCode int, m *gomol.Attrs, msg string, args ...interface{}) {
	a.Log(gomol.LevelFatal, m, msg, args...)
	a.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
