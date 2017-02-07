package gomolreplay

import (
	"time"

	"github.com/aphistic/gomol"

	. "gopkg.in/check.v1"
)

type (
	logFunc  func(msg string) error
	logFuncf func(msg string, a ...interface{}) error
	logFuncm func(attrs *gomol.Attrs, msg string, a ...interface{}) error
)

var (
	AllLevels = []gomol.LogLevel{
		gomol.LevelDebug,
		gomol.LevelInfo,
		gomol.LevelWarning,
		gomol.LevelError,
		gomol.LevelFatal,
	}
)

func (s *ReplaySuite) TestConvenienceMethods(c *C) {
	var (
		logger   = newDefaultMockLogger()
		replay   = NewReplayAdapter(logger, AllLevels...)
		messages = []logArgs{}
	)

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, logArgs{level, attrs, msg, a})
		return nil
	}

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, logArgs{level, attrs, msg, a})
		return nil
	}

	for level, method := range map[gomol.LogLevel]logFunc{
		gomol.LevelDebug:   replay.Dbg,
		gomol.LevelInfo:    replay.Info,
		gomol.LevelWarning: replay.Warn,
		gomol.LevelError:   replay.Err,
		gomol.LevelFatal:   replay.Fatal,
	} {
		method("foo")
		method("bar")
		method("baz")
		replay.Replay(gomol.LevelFatal)
		assert(c, messages, level, nil, nil)

		// Reset
		replay.reset()
		messages = messages[:0]
	}

	for level, method := range map[gomol.LogLevel]logFuncf{
		gomol.LevelDebug:   replay.Dbgf,
		gomol.LevelInfo:    replay.Infof,
		gomol.LevelWarning: replay.Warnf,
		gomol.LevelError:   replay.Errf,
		gomol.LevelFatal:   replay.Fatalf,
	} {
		method("foo", 12)
		method("bar", 43)
		method("baz", 74)
		replay.Replay(gomol.LevelFatal)
		assert(c, messages, level, nil, []int{12, 43, 74})

		// Reset
		replay.reset()
		messages = messages[:0]
	}

	for level, method := range map[gomol.LogLevel]logFuncm{
		gomol.LevelDebug:   replay.Dbgm,
		gomol.LevelInfo:    replay.Infom,
		gomol.LevelWarning: replay.Warnm,
		gomol.LevelError:   replay.Errm,
		gomol.LevelFatal:   replay.Fatalm,
	} {
		method(gomol.NewAttrsFromMap(map[string]interface{}{"x": "x"}), "foo", 12)
		method(gomol.NewAttrsFromMap(map[string]interface{}{"y": "y"}), "bar", 43)
		method(gomol.NewAttrsFromMap(map[string]interface{}{"z": "z"}), "baz", 74)
		replay.Replay(gomol.LevelFatal)
		assert(c, messages, level, []string{"x", "y", "z"}, []int{12, 43, 74})

		// Reset
		replay.reset()
		messages = messages[:0]
	}
}

func assert(c *C, messages []logArgs, level gomol.LogLevel, attrs []string, params []int) {
	for i := 0; i < 3; i++ {
		c.Assert(messages[i+0].level, Equals, level)
		c.Assert(messages[i+3].level, Equals, gomol.LevelFatal)
	}

	for i, msg := range []string{"foo", "bar", "baz"} {
		c.Assert(messages[i+0].msg, Equals, msg)
		c.Assert(messages[i+3].msg, Equals, msg)
	}

	for i, val := range params {
		c.Assert(messages[i+0].a[0], Equals, val)
		c.Assert(messages[i+3].a[0], Equals, val)
	}

	for i, val := range attrs {
		c.Assert(messages[i+0].attrs.GetAttr(val), Equals, val)
		c.Assert(messages[i+3].attrs.GetAttr(val), Equals, val)
	}
}
