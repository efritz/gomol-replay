package gomolreplay

import (
	"time"

	"github.com/aphistic/gomol"
)

const (
	// AttrReplay is an attribute assigned to a message that has
	// been replayed at a different log level. Its value is equal
	// to the original log level.
	AttrReplay = "replayed-from-level"
)

type (
	// Adapter provides a way to replay a sequence of message, in
	// the order they were logged, at a higher log level.
	Adapter struct {
		base            gomol.WrappableLogger
		clock           clock
		journal         []*logMessage
		journaledLevels []gomol.LogLevel
		replayingAt     *gomol.LogLevel
	}

	logMessage struct {
		level gomol.LogLevel
		attrs *gomol.Attrs
		ts    time.Time
		msg   string
		args  []interface{}
	}
)

// NewAdapter creates an Adapter which wraps the given logger.
func NewAdapter(logger gomol.WrappableLogger, journaledLevels ...gomol.LogLevel) *Adapter {
	return newAdapterWithClock(logger, &realClock{}, journaledLevels...)
}

func newAdapterWithClock(logger gomol.WrappableLogger, clock clock, journaledLevels ...gomol.LogLevel) *Adapter {
	return &Adapter{
		base:            logger,
		clock:           clock,
		journal:         []*logMessage{},
		journaledLevels: journaledLevels,
	}
}

// LogWithTime will log a message at the provided level to all loggers added
// to the logger wrapped by this RollupAdapter. It is similar to Log except
// the timestamp will be set to the value of ts.
func (a *Adapter) LogWithTime(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, args ...interface{}) error {
	if err := a.base.LogWithTime(level, ts, attrs, msg, args...); err != nil {
		return err
	}

	if a.shouldJournal(level) {
		message := &logMessage{level: level, attrs: attrs, ts: ts, msg: msg, args: args}

		if a.replayingAt != nil {
			if err := a.replayMessage(message); err != nil {
				return err
			}
		}

		a.journal = append(a.journal, message)
	}

	return nil
}

// Log will log a message at the provided level to all loggers added to the
// logger wrapped by this RollupAdapter.
func (a *Adapter) Log(level gomol.LogLevel, attrs *gomol.Attrs, msg string, args ...interface{}) error {
	if !a.shouldJournal(level) {
		return a.base.Log(level, attrs, msg, args...)
	}

	return a.LogWithTime(level, a.clock.Now(), attrs, msg, args...)
}

// ShutdownLoggers will call the wrapped logger's ShutdownLoggers method.
func (a *Adapter) ShutdownLoggers() error {
	return a.base.ShutdownLoggers()
}

// Replay will cause all of the messages previously logged at one of the
// journaled levels to be re-set at the given level. All future messages
// logged at one of the journaled levels will be replayed immediately.
func (a *Adapter) Replay(level gomol.LogLevel) error {
	if a.replayingAt != nil && *a.replayingAt <= level {
		return nil
	}

	a.replayingAt = &level

	for _, message := range a.journal {
		if err := a.replayMessage(message); err != nil {
			return err
		}
	}

	return nil
}

func (a *Adapter) shouldJournal(level gomol.LogLevel) bool {
	for _, l := range a.journaledLevels {
		if l == level {
			return true
		}
	}

	return false
}

func (a *Adapter) replayMessage(message *logMessage) error {
	return a.base.LogWithTime(*a.replayingAt, message.ts, addAttr(message.attrs, message.level), message.msg, message.args...)
}

func addAttr(attrs *gomol.Attrs, level gomol.LogLevel) *gomol.Attrs {
	if attrs == nil {
		attrs = gomol.NewAttrs()
	} else {
		attrs = gomol.NewAttrsFromAttrs(attrs)
	}

	return attrs.SetAttr(AttrReplay, level)
}

func (a *Adapter) reset() {
	a.journal = a.journal[:0]
	a.replayingAt = nil
}
