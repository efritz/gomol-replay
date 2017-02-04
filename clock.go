package gomolreplay

import "time"

type (
	clock interface {
		Now() time.Time
	}

	realClock struct{}
)

func (rc *realClock) Now() time.Time {
	return time.Now()
}
