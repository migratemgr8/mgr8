package infrastructure

import "time"

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func NewClock() *realClock {
	return &realClock{}
}

func (*realClock) Now() time.Time {
	return time.Now()
}
