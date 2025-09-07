package limit

import "time"

type sysClock struct {
}

func NewClock() sysClock {
	return sysClock{}
}

func (c sysClock) now() time.Time {
	return time.Now()
}
