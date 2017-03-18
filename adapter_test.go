package gomolreplay

import (
	"fmt"
	"time"

	"github.com/aphistic/gomol"

	. "gopkg.in/check.v1"
)

func (s *ReplaySuite) TestWhitelistLevelsAreJournaled(c *C) {
	var (
		logger   = newDefaultMockLogger()
		adapter  = NewAdapter(logger, gomol.LevelInfo, gomol.LevelError)
		messages = []string{}
	)

	logger.log = func(level gomol.LogLevel, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, msg)
		return nil
	}

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, msg)
		return nil
	}

	adapter.Log(gomol.LevelDebug, nil, "foo")
	adapter.Log(gomol.LevelInfo, nil, "bar")
	adapter.Log(gomol.LevelWarning, nil, "baz")
	adapter.Log(gomol.LevelError, nil, "bnk")
	adapter.Log(gomol.LevelFatal, nil, "qux")

	c.Assert(len(messages), Equals, 5)
	c.Assert(messages[0], Equals, "foo")
	c.Assert(messages[1], Equals, "bar")
	c.Assert(messages[2], Equals, "baz")
	c.Assert(messages[3], Equals, "bnk")
	c.Assert(messages[4], Equals, "qux")

	c.Assert(len(adapter.journal), Equals, 2)
	c.Assert(adapter.journal[0].msg, Equals, "bar")
	c.Assert(adapter.journal[1].msg, Equals, "bnk")
}

func (s *ReplaySuite) TestReplayJournal(c *C) {
	var (
		logger   = newDefaultMockLogger()
		adapter  = NewAdapter(logger, gomol.LevelDebug)
		messages = []logArgs{}
	)

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, logArgs{level, attrs, msg, a})
		return nil
	}

	adapter.Log(gomol.LevelDebug, gomol.NewAttrsFromMap(map[string]interface{}{"x": "x"}), "foo", 12)
	adapter.Log(gomol.LevelDebug, gomol.NewAttrsFromMap(map[string]interface{}{"y": "y"}), "bar", 43)
	adapter.Log(gomol.LevelDebug, gomol.NewAttrsFromMap(map[string]interface{}{"z": "z"}), "baz", 74)
	adapter.Replay(gomol.LevelWarning)

	c.Assert(len(messages), Equals, 6)

	for i := 0; i < 3; i++ {
		c.Assert(messages[i+0].level, Equals, gomol.LevelDebug)
		c.Assert(messages[i+3].level, Equals, gomol.LevelWarning)
	}

	for i, msg := range []string{"foo", "bar", "baz"} {
		c.Assert(messages[i+0].msg, Equals, msg)
		c.Assert(messages[i+3].msg, Equals, msg)
	}

	for i, val := range []int{12, 43, 74} {
		c.Assert(messages[i+0].a[0], Equals, val)
		c.Assert(messages[i+3].a[0], Equals, val)
	}

	for i, val := range []string{"x", "y", "z"} {
		c.Assert(messages[i+0].attrs.GetAttr(val), Equals, val)
		c.Assert(messages[i+3].attrs.GetAttr(val), Equals, val)
	}
}

func (s *ReplaySuite) TestReplayTwice(c *C) {
	var (
		logger   = newDefaultMockLogger()
		adapter  = NewAdapter(logger, gomol.LevelDebug)
		messages = []logArgs{}
	)

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, logArgs{level, attrs, msg, a})
		return nil
	}

	adapter.Log(gomol.LevelDebug, nil, "foo")
	adapter.Log(gomol.LevelDebug, nil, "bar")
	adapter.Log(gomol.LevelDebug, nil, "baz")
	adapter.Replay(gomol.LevelWarning)
	adapter.Replay(gomol.LevelError)

	c.Assert(len(messages), Equals, 9)
	c.Assert(messages[0].level, Equals, gomol.LevelDebug)
	c.Assert(messages[1].level, Equals, gomol.LevelDebug)
	c.Assert(messages[2].level, Equals, gomol.LevelDebug)
	c.Assert(messages[3].level, Equals, gomol.LevelWarning)
	c.Assert(messages[4].level, Equals, gomol.LevelWarning)
	c.Assert(messages[5].level, Equals, gomol.LevelWarning)
	c.Assert(messages[6].level, Equals, gomol.LevelError)
	c.Assert(messages[7].level, Equals, gomol.LevelError)
	c.Assert(messages[8].level, Equals, gomol.LevelError)

	for i, msg := range []string{"foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz"} {
		c.Assert(messages[i].msg, Equals, msg)
	}
}

func (s *ReplaySuite) TestReplayAtHigherlevelNoops(c *C) {
	var (
		logger   = newDefaultMockLogger()
		adapter  = NewAdapter(logger, gomol.LevelDebug)
		messages = []logArgs{}
	)

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, logArgs{level, attrs, msg, a})
		return nil
	}

	adapter.Log(gomol.LevelDebug, nil, "foo")
	adapter.Log(gomol.LevelDebug, nil, "bar")
	adapter.Log(gomol.LevelDebug, nil, "baz")
	adapter.Replay(gomol.LevelError)
	adapter.Replay(gomol.LevelWarning)

	c.Assert(len(messages), Equals, 6)
	c.Assert(messages[0].level, Equals, gomol.LevelDebug)
	c.Assert(messages[1].level, Equals, gomol.LevelDebug)
	c.Assert(messages[2].level, Equals, gomol.LevelDebug)
	c.Assert(messages[3].level, Equals, gomol.LevelError)
	c.Assert(messages[4].level, Equals, gomol.LevelError)
	c.Assert(messages[5].level, Equals, gomol.LevelError)

	for i, msg := range []string{"foo", "bar", "baz", "foo", "bar", "baz"} {
		c.Assert(messages[i].msg, Equals, msg)
	}
}

func (s *ReplaySuite) TestLogAfterReplaySendsImmediately(c *C) {
	var (
		logger   = newDefaultMockLogger()
		adapter  = NewAdapter(logger, gomol.LevelDebug)
		messages = []logArgs{}
	)

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, logArgs{level, attrs, msg, a})
		return nil
	}

	adapter.Log(gomol.LevelDebug, nil, "foo")
	adapter.Log(gomol.LevelDebug, nil, "bar")
	adapter.Log(gomol.LevelDebug, nil, "baz")
	adapter.Replay(gomol.LevelWarning)
	adapter.Log(gomol.LevelDebug, nil, "bnk")
	adapter.Log(gomol.LevelDebug, nil, "qux")

	c.Assert(len(messages), Equals, 10)
	c.Assert(messages[0].level, Equals, gomol.LevelDebug)
	c.Assert(messages[1].level, Equals, gomol.LevelDebug)
	c.Assert(messages[2].level, Equals, gomol.LevelDebug)
	c.Assert(messages[3].level, Equals, gomol.LevelWarning)
	c.Assert(messages[4].level, Equals, gomol.LevelWarning)
	c.Assert(messages[5].level, Equals, gomol.LevelWarning)
	c.Assert(messages[6].level, Equals, gomol.LevelDebug)
	c.Assert(messages[7].level, Equals, gomol.LevelWarning)
	c.Assert(messages[8].level, Equals, gomol.LevelDebug)
	c.Assert(messages[9].level, Equals, gomol.LevelWarning)

	for i, msg := range []string{"foo", "bar", "baz", "foo", "bar", "baz", "bnk", "bnk", "qux", "qux"} {
		c.Assert(messages[i].msg, Equals, msg)
	}
}

func (s *ReplaySuite) TestLogAfterSecondReplaySendsAtNewLevel(c *C) {
	var (
		logger   = newDefaultMockLogger()
		adapter  = NewAdapter(logger, gomol.LevelDebug)
		messages = []logArgs{}
	)

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, logArgs{level, attrs, msg, a})
		return nil
	}

	adapter.Log(gomol.LevelDebug, nil, "foo")
	adapter.Log(gomol.LevelDebug, nil, "bar")
	adapter.Replay(gomol.LevelWarning)
	adapter.Replay(gomol.LevelError)
	adapter.Log(gomol.LevelDebug, nil, "baz")
	adapter.Log(gomol.LevelDebug, nil, "bnk")

	c.Assert(len(messages), Equals, 10)
	c.Assert(messages[0].level, Equals, gomol.LevelDebug)
	c.Assert(messages[1].level, Equals, gomol.LevelDebug)
	c.Assert(messages[2].level, Equals, gomol.LevelWarning)
	c.Assert(messages[3].level, Equals, gomol.LevelWarning)
	c.Assert(messages[4].level, Equals, gomol.LevelError)
	c.Assert(messages[5].level, Equals, gomol.LevelError)
	c.Assert(messages[6].level, Equals, gomol.LevelDebug)
	c.Assert(messages[7].level, Equals, gomol.LevelError)
	c.Assert(messages[8].level, Equals, gomol.LevelDebug)
	c.Assert(messages[9].level, Equals, gomol.LevelError)

	for i, msg := range []string{"foo", "bar", "foo", "bar", "foo", "bar", "baz", "baz", "bnk", "bnk"} {
		c.Assert(messages[i].msg, Equals, msg)
	}
}

func (s *ReplaySuite) TestReplayKeepsOriginalTimestamp(c *C) {
	var (
		logger  = newDefaultMockLogger()
		clock   = newMockClock(24000)
		adapter = newAdapterWithClock(logger, clock, gomol.LevelDebug)
		times1  = map[string]time.Time{}
		times2  = map[string]time.Time{}
	)

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		if level == gomol.LevelDebug {
			times1[msg] = ts
		} else {
			times2[msg] = ts
		}

		return nil
	}

	adapter.Log(gomol.LevelDebug, nil, "foo")
	adapter.LogWithTime(gomol.LevelDebug, time.Unix(61, 500), nil, "bar")
	clock.advance(3000)
	adapter.Log(gomol.LevelDebug, nil, "baz")
	clock.advance(10000)
	adapter.Replay(gomol.LevelError)

	c.Assert(len(times1), Equals, 3)
	c.Assert(len(times2), Equals, 3)
	c.Assert(times1["foo"], Equals, time.Unix(24, 0))
	c.Assert(times2["foo"], Equals, time.Unix(24, 0))
	c.Assert(times1["bar"], Equals, time.Unix(61, 500))
	c.Assert(times2["bar"], Equals, time.Unix(61, 500))
	c.Assert(times1["baz"], Equals, time.Unix(27, 0))
	c.Assert(times2["baz"], Equals, time.Unix(27, 0))
}

func (s *ReplaySuite) TestCheckReplayAddsAttribute(c *C) {
	var (
		logger   = newDefaultMockLogger()
		adapter  = NewAdapter(logger, gomol.LevelDebug, gomol.LevelInfo)
		messages = []logArgs{}
	)

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, logArgs{level, attrs, msg, a})
		return nil
	}

	adapter.Log(gomol.LevelDebug, nil, "foo")
	adapter.Log(gomol.LevelInfo, nil, "bar")
	adapter.Log(gomol.LevelDebug, nil, "baz")
	adapter.Replay(gomol.LevelError)
	adapter.Log(gomol.LevelDebug, nil, "bnk")

	c.Assert(len(messages), Equals, 8)
	c.Assert(messages[0].attrs, IsNil)
	c.Assert(messages[1].attrs, IsNil)
	c.Assert(messages[2].attrs, IsNil)
	c.Assert(messages[3].attrs.GetAttr(AttrReplay), Equals, gomol.LevelDebug)
	c.Assert(messages[4].attrs.GetAttr(AttrReplay), Equals, gomol.LevelInfo)
	c.Assert(messages[5].attrs.GetAttr(AttrReplay), Equals, gomol.LevelDebug)
	c.Assert(messages[6].attrs, IsNil)
	c.Assert(messages[7].attrs.GetAttr(AttrReplay), Equals, gomol.LevelDebug)
}

func (s *ReplaySuite) TestCheckSecondReplayAddsAttribute(c *C) {
	var (
		logger   = newDefaultMockLogger()
		adapter  = NewAdapter(logger, gomol.LevelDebug, gomol.LevelInfo)
		messages = []logArgs{}
	)

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		messages = append(messages, logArgs{level, attrs, msg, a})
		return nil
	}

	adapter.Log(gomol.LevelDebug, nil, "foo")
	adapter.Log(gomol.LevelInfo, nil, "bar")
	adapter.Replay(gomol.LevelWarning)
	adapter.Replay(gomol.LevelError)
	adapter.Log(gomol.LevelDebug, nil, "bnk")

	c.Assert(len(messages), Equals, 8)
	c.Assert(messages[0].attrs, IsNil)
	c.Assert(messages[1].attrs, IsNil)
	c.Assert(messages[2].attrs.GetAttr(AttrReplay), Equals, gomol.LevelDebug)
	c.Assert(messages[3].attrs.GetAttr(AttrReplay), Equals, gomol.LevelInfo)
	c.Assert(messages[4].attrs.GetAttr(AttrReplay), Equals, gomol.LevelDebug)
	c.Assert(messages[5].attrs.GetAttr(AttrReplay), Equals, gomol.LevelInfo)
	c.Assert(messages[6].attrs, IsNil)
	c.Assert(messages[7].attrs.GetAttr(AttrReplay), Equals, gomol.LevelDebug)
}

func (s *ReplaySuite) TestReplayErrors(c *C) {
	var (
		logger  = newDefaultMockLogger()
		adapter = NewAdapter(logger, gomol.LevelDebug, gomol.LevelInfo)
		calls   = 0
	)

	logger.logWithTime = func(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
		calls++
		if calls > 3 {
			return fmt.Errorf("Error %d", calls-3)
		}

		return nil
	}

	c.Assert(adapter.LogWithTime(gomol.LevelInfo, time.Now(), nil, "foo"), IsNil)
	c.Assert(adapter.LogWithTime(gomol.LevelInfo, time.Now(), nil, "bar"), IsNil)
	c.Assert(adapter.LogWithTime(gomol.LevelInfo, time.Now(), nil, "baz"), IsNil)

	c.Assert(adapter.Replay(gomol.LevelError), ErrorMatches, "Error 1")
	c.Assert(adapter.LogWithTime(gomol.LevelInfo, time.Now(), nil, "baz"), ErrorMatches, "Error 2")
}

func (s *ReplaySuite) TestShutdownLoggers(c *C) {
	var (
		logger  = newDefaultMockLogger()
		adapter = NewAdapter(logger)
	)

	logger.shutdownLoggers = func() error {
		return fmt.Errorf("foo")
	}

	c.Assert(adapter.ShutdownLoggers(), ErrorMatches, "foo")
}
