package gomolreplay

import (
	"testing"
	"time"

	"github.com/aphistic/gomol"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ReplaySuite struct{}

var _ = Suite(&ReplaySuite{})

//
//

type logArgs struct {
	level gomol.LogLevel
	attrs *gomol.Attrs
	msg   string
	a     []interface{}
}

//
// Mocks

type mockLogger struct {
	logWithTime     func(gomol.LogLevel, time.Time, *gomol.Attrs, string, ...interface{}) error
	log             func(gomol.LogLevel, *gomol.Attrs, string, ...interface{}) error
	shutdownLoggers func() error
}

func newDefaultMockLogger() *mockLogger {
	return &mockLogger{
		logWithTime:     func(gomol.LogLevel, time.Time, *gomol.Attrs, string, ...interface{}) error { return nil },
		log:             func(gomol.LogLevel, *gomol.Attrs, string, ...interface{}) error { return nil },
		shutdownLoggers: func() error { return nil },
	}
}

func (m *mockLogger) LogWithTime(level gomol.LogLevel, ts time.Time, attrs *gomol.Attrs, msg string, a ...interface{}) error {
	return m.logWithTime(level, ts, attrs, msg, a...)
}

func (m *mockLogger) Log(level gomol.LogLevel, attrs *gomol.Attrs, msg string, a ...interface{}) error {
	return m.log(level, attrs, msg, a...)
}

func (m *mockLogger) ShutdownLoggers() error {
	return m.shutdownLoggers()
}

type mockClock struct {
	seconds     int64
	nanoseconds int64
}

func newMockClock(ms int64) *mockClock {
	m := &mockClock{}
	m.advance(ms)
	return m
}

func (m *mockClock) advance(ms int64) {
	m.seconds += ms / 1000
	m.nanoseconds += (ms % 1000) * 1e6

	for m.nanoseconds >= 1e9 {
		m.seconds++
		m.nanoseconds -= 1e9
	}
}

func (m *mockClock) Now() time.Time {
	return time.Unix(m.seconds, m.nanoseconds)
}
