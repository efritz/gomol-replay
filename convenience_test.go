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

	testExiter struct {
		exited bool
		code   int
	}

	dieTestFunc struct {
		f           func()
		checkAttrs  bool
		checkParams bool
	}
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

func (exiter *testExiter) Exit(code int) {
	exiter.code = code
	exiter.exited = true
}

func (s *ReplaySuite) TestDieEndsProcessAfterImmediateReplay(c *C) {
	var (
		messages  = []logArgs{}
		exitCount = 0

		exiter = &testExiter{}
		logger = newDefaultMockLogger()
		replay = NewReplayAdapter(logger, AllLevels...)
	)

	setExiter(exiter)

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, logArgs{level, attrs, msg, a})
		return nil
	}

	logger.shutdownLoggers = func() error {
		exitCount++
		return nil
	}

	for i, data := range []dieTestFunc{
		{func() { replay.Die(1000, "foo") }, false, false},
		{func() { replay.Dief(2000, "foo", 42) }, false, true},
		{func() { replay.Diem(3000, gomol.NewAttrsFromMap(map[string]interface{}{"x": "y"}), "foo", 42) }, true, false},
	} {
		replay.Replay(gomol.LevelWarning)
		data.f()

		c.Assert(exitCount, Equals, i+1)
		c.Assert(exiter.code, Equals, 1000*(i+1))

		c.Assert(len(messages), Equals, 2)
		c.Assert(messages[0].level, Equals, gomol.LevelFatal)
		c.Assert(messages[1].level, Equals, gomol.LevelWarning)
		c.Assert(messages[0].msg, Equals, "foo")
		c.Assert(messages[1].msg, Equals, "foo")

		if data.checkAttrs {
			c.Assert(messages[0].attrs.GetAttr("x"), Equals, "y")
			c.Assert(messages[1].attrs.GetAttr("x"), Equals, "y")
		}

		if data.checkParams {
			c.Assert(messages[0].a[0], Equals, 42)
			c.Assert(messages[1].a[0], Equals, 42)
		}

		// Reset
		messages = messages[:0]
	}
}

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
	c.Assert(len(messages), Equals, 6)

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
