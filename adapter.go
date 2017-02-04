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
	// ReplayAdapter provides a way to replay a sequence of message, in
	// the order they were logged, at a higher log level.
	ReplayAdapter struct {
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
		a     []interface{}
	}
)

// NewReplayAdapter creates a ReplayAdapter which wraps the given logger.
func NewReplayAdapter(logger gomol.WrappableLogger, journaledLevels ...gomol.LogLevel) *ReplayAdapter {
	return newReplayAdapterWithClock(logger, &realClock{}, journaledLevels...)
}

func newReplayAdapterWithClock(logger gomol.WrappableLogger, clock clock, journaledLevels ...gomol.LogLevel) *ReplayAdapter {
	return &ReplayAdapter{
		base:            logger,
		clock:           clock,
		journal:         []*logMessage{},
		journaledLevels: journaledLevels,
	}
}

func (ra *ReplayAdapter) LogWithTime(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
	if err := ra.base.LogWithTime(level, ts, attrs, msg, a...); err != nil {
		return err
	}

	if ra.shouldJournal(level) {
		message := &logMessage{level: level, attrs: attrs, ts: ts, msg: msg, a: a}

		if ra.replayingAt != nil {
			if err := ra.replayMessage(message); err != nil {
				return err
			}
		}

		ra.journal = append(ra.journal, message)
	}

	return nil
}

func (ra *ReplayAdapter) Log(level gomol.LogLevel, attrs *gomol.Attrs, msg string, a ...interface{}) error {
	if !ra.shouldJournal(level) {
		return ra.base.Log(level, attrs, msg, a...)
	}

	return ra.LogWithTime(level, ra.clock.Now(), attrs, msg, a...)
}

func (ra *ReplayAdapter) ShutdownLoggers() error {
	return ra.base.ShutdownLoggers()
}

func (ra *ReplayAdapter) Replay(level gomol.LogLevel) error {
	if ra.replayingAt != nil && *ra.replayingAt <= level {
		return nil
	}

	ra.replayingAt = &level

	for _, message := range ra.journal {
		if err := ra.replayMessage(message); err != nil {
			return err
		}
	}

	return nil
}

func (ra *ReplayAdapter) shouldJournal(level gomol.LogLevel) bool {
	for _, l := range ra.journaledLevels {
		if l == level {
			return true
		}
	}

	return false
}

func (ra *ReplayAdapter) replayMessage(message *logMessage) error {
	return ra.base.LogWithTime(*ra.replayingAt, message.ts, addAttr(message.attrs, message.level), message.msg, message.a...)
}

func addAttr(attrs *gomol.Attrs, level gomol.LogLevel) *gomol.Attrs {
	if attrs == nil {
		attrs = gomol.NewAttrs()
	} else {
		attrs = gomol.NewAttrsFromAttrs(attrs)
	}

	return attrs.SetAttr(AttrReplay, level)
}
