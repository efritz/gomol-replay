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
func (ra *ReplayAdapter) Dbg(msg string) error {
	return ra.Debug(msg)
}

// Dbgf is a short-hand version of Debugf
func (ra *ReplayAdapter) Dbgf(msg string, a ...interface{}) error {
	return ra.Debugf(msg, a...)
}

// Dbgm is a short-hand version of Debugm
func (ra *ReplayAdapter) Dbgm(m *gomol.Attrs, msg string, a ...interface{}) error {
	return ra.Debugm(m, msg, a...)
}

// Debug logs msg to all added loggers at LogLevel.LevelDebug
func (ra *ReplayAdapter) Debug(msg string) error {
	return ra.Log(gomol.LevelDebug, nil, msg)
}

/*
Debugf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelDebug
*/
func (ra *ReplayAdapter) Debugf(msg string, a ...interface{}) error {
	return ra.Log(gomol.LevelDebug, nil, msg, a...)
}

/*
Debugm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelDebug. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (ra *ReplayAdapter) Debugm(m *gomol.Attrs, msg string, a ...interface{}) error {
	return ra.Log(gomol.LevelDebug, m, msg, a...)
}

// Info logs msg to all added loggers at LogLevel.LevelInfo
func (ra *ReplayAdapter) Info(msg string) error {
	return ra.Log(gomol.LevelInfo, nil, msg)
}

/*
Infof uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelInfo
*/
func (ra *ReplayAdapter) Infof(msg string, a ...interface{}) error {
	return ra.Log(gomol.LevelInfo, nil, msg, a...)
}

/*
Infom uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelInfo. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (ra *ReplayAdapter) Infom(m *gomol.Attrs, msg string, a ...interface{}) error {
	return ra.Log(gomol.LevelInfo, m, msg, a...)
}

// Warn is a short-hand version of Warning
func (ra *ReplayAdapter) Warn(msg string) error {
	return ra.Warning(msg)
}

// Warnf is a short-hand version of Warningf
func (ra *ReplayAdapter) Warnf(msg string, a ...interface{}) error {
	return ra.Warningf(msg, a...)
}

// Warnm is a short-hand version of Warningm
func (ra *ReplayAdapter) Warnm(m *gomol.Attrs, msg string, a ...interface{}) error {
	return ra.Warningm(m, msg, a...)
}

/*
Warning uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning
*/
func (ra *ReplayAdapter) Warning(msg string) error {
	return ra.Log(gomol.LevelWarning, nil, msg)
}

/*
Warningf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning
*/
func (ra *ReplayAdapter) Warningf(msg string, a ...interface{}) error {
	return ra.Log(gomol.LevelWarning, nil, msg, a...)
}

/*
Warningm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (ra *ReplayAdapter) Warningm(m *gomol.Attrs, msg string, a ...interface{}) error {
	return ra.Log(gomol.LevelWarning, m, msg, a...)
}

// Err is a short-hand version of Error
func (ra *ReplayAdapter) Err(msg string) error {
	return ra.Error(msg)
}

// Errf is a short-hand version of Errorf
func (ra *ReplayAdapter) Errf(msg string, a ...interface{}) error {
	return ra.Errorf(msg, a...)
}

// Errm is a short-hand version of Errorm
func (ra *ReplayAdapter) Errm(m *gomol.Attrs, msg string, a ...interface{}) error {
	return ra.Errorm(m, msg, a...)
}

/*
Error uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError
*/
func (ra *ReplayAdapter) Error(msg string) error {
	return ra.Log(gomol.LevelError, nil, msg)
}

/*
Errorf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError
*/
func (ra *ReplayAdapter) Errorf(msg string, a ...interface{}) error {
	return ra.Log(gomol.LevelError, nil, msg, a...)
}

/*
Errorm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (ra *ReplayAdapter) Errorm(m *gomol.Attrs, msg string, a ...interface{}) error {
	return ra.Log(gomol.LevelError, m, msg, a...)
}

/*
Fatal uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal
*/
func (ra *ReplayAdapter) Fatal(msg string) error {
	return ra.Log(gomol.LevelFatal, nil, msg)
}

/*
Fatalf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal
*/
func (ra *ReplayAdapter) Fatalf(msg string, a ...interface{}) error {
	return ra.Log(gomol.LevelFatal, nil, msg, a...)
}

/*
Fatalm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (ra *ReplayAdapter) Fatalm(m *gomol.Attrs, msg string, a ...interface{}) error {
	return ra.Log(gomol.LevelFatal, m, msg, a...)
}

// Die will log a message using Fatal, call ShutdownLoggers and then exit the application with the provided exit code.
// This function is not subject to rollup and is always sent to the wrapped logger.
func (ra *ReplayAdapter) Die(exitCode int, msg string) {
	ra.Log(gomol.LevelFatal, nil, msg)
	ra.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}

// Dief will log a message using Fatalf, call ShutdownLoggers and then exit the application with the provided exit code.
func (ra *ReplayAdapter) Dief(exitCode int, msg string, a ...interface{}) {
	ra.Log(gomol.LevelFatal, nil, msg, a...)
	ra.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}

// Diem will log a message using Fatalm, call ShutdownLoggers and then exit the application with the provided exit code.
func (ra *ReplayAdapter) Diem(exitCode int, m *gomol.Attrs, msg string, a ...interface{}) {
	ra.Log(gomol.LevelFatal, m, msg, a...)
	ra.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
