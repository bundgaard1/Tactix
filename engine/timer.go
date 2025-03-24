package engine

import "time"

const (
	InfiniteTime = -1
	NoValue      = 0
)

type Timer struct {
	// UCI options, theses are apparently in nanoseconds
	TimeLeft  int64
	Increment int64
	MoveTime  int64
	MovesToGo int64
	MaxDepth  uint8

	Stop        bool
	TimeForMove int64
	stopTime    time.Time
}

func NewTimer() (tm Timer) {
	return tm
}

func (tm *Timer) NoTimeControl() bool {
	return tm.TimeLeft == InfiniteTime && tm.MoveTime == NoValue
}

func (tm *Timer) SetTimeForMove(timeForMove int64) {
	tm.stopTime = time.Now().Add(time.Duration(timeForMove) * time.Millisecond)
	tm.TimeForMove = timeForMove

}

func (tm *Timer) SetTimeControl(timeLeft, increment, movetime, movesToGo int64, maxDepth uint8) {
	tm.TimeLeft = timeLeft
	tm.Increment = increment
	tm.MoveTime = movetime
	tm.MovesToGo = movesToGo
	tm.MaxDepth = maxDepth
}

func (tm *Timer) Start() {
	tm.Stop = false

	if tm.MoveTime != NoValue {
		var buffer int64 = 100
		tm.stopTime = time.Now().Add(time.Duration(tm.MoveTime-buffer) * time.Millisecond)
		return
	}

	if tm.TimeLeft == InfiniteTime {
		return
	}
	timeForMove := 0
	if tm.MovesToGo != NoValue {
		timeForMove = int(tm.TimeLeft / tm.MovesToGo)
	} else {
		timeForMove = int(tm.TimeLeft / 30)
	}
	if tm.Increment != NoValue {
		timeForMove += int(tm.Increment)
	}

	tm.stopTime = time.Now().Add(time.Duration(timeForMove) * time.Millisecond)
	tm.TimeForMove = int64(timeForMove)

}

func (tm *Timer) Check() {
	if tm.MoveTime == NoValue {
		return
	}
	if time.Now().After(tm.stopTime) {
		tm.Stop = true
	}
}
