package ratelimiter

import "time"

type Clock interface {
	Now() time.Time
}

type defaultClock struct {
}

type mockClock struct {
	currentTime time.Time
}

func NewClock() Clock {
	return defaultClock{}
}

func NewMockClock() Clock {
	return &mockClock{
		currentTime: time.Now(),
	}
}

func (c defaultClock) Now() time.Time {
	return time.Now()
}

func (c mockClock) Now() time.Time {
	return time.Now()
}

func (m *mockClock) Advance(d time.Duration) {
	m.currentTime = m.currentTime.Add(d)
}
