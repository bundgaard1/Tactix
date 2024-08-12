package engine

import "time"

type Timer struct {
	// UCI options, theses are apparently in nanoseconds
	TimeLeft int64
	MoveTime int64

	Stop        bool
	TimeForMove int64
	stopTime    time.Time
}

func NewTimer() (tm Timer) {
	return tm
}

func (tm *Timer) Start() {
	tm.Stop = false

	tm.stopTime = time.Now().Add(time.Duration(tm.MoveTime) * time.Millisecond)
}

func (tm *Timer) Check() {
	if time.Now().After(tm.stopTime) {
		tm.Stop = true
	}
}
