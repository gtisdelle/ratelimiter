package ratelimiter

import "time"

type Clock interface {
	Now() time.Time
}

type defaultClock struct {
}

type MockClock struct {
	currentTime time.Time
}

func NewClock() Clock {
	return defaultClock{}
}

func (c defaultClock) Now() time.Time {
	return time.Now()
}

func (c *MockClock) Now() time.Time {
	return c.currentTime
}

func (m *MockClock) Advance(d time.Duration) {
	m.currentTime = m.currentTime.Add(d)
}
