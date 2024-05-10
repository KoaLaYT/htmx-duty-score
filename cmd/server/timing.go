package main

import (
	"fmt"
	"time"
)

type TimingType string

const (
	CONSUME TimingType = "consume"
	TOPUP   TimingType = "topup"
)

func (t TimingType) String() string {
	return string(t)
}

type Timing struct {
	Start    time.Time
	Stop     time.Time
	TimeType TimingType
}

func (t Timing) IsStopped() bool {
	return !t.Stop.IsZero()
}

func (t Timing) StartString() string {
	return t.Start.Format("15:04")
}

func (t Timing) StopString() string {
	if t.Stop.IsZero() {
		return ""
	}
	return t.Start.Format("15:04")
}

func (t Timing) Duration() string {
	if t.Stop.IsZero() {
		return ""
	}
	duration := t.Stop.Sub(t.Start).Truncate(time.Second).String()
	unit := "+"
	if t.TimeType == CONSUME {
		unit = "-"
	}
	return fmt.Sprintf("%s%s", unit, duration)
}
