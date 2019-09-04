package main

import "C"

import (
	"fmt"
	"time"
)

type IntervalTimer struct {
	Interval time.Duration
	Enabled bool
	Job 	func()
}

func (it *IntervalTimer) Isr() {
	if it.Enabled {
		it.Job()
		time.AfterFunc(it.Interval, it.Isr)
	}
}

func (it *IntervalTimer) Start() {
	it.Enabled = true
	time.AfterFunc(it.Interval, it.Isr)
}

func (it *IntervalTimer) Stop() {
	it.Enabled = false
}

var it *IntervalTimer

func PrintTask() {
	fmt.Println("go timer is on")
}

//export StartTimer
func StartTimer(millis int) {
	it = &IntervalTimer{Interval: time.Duration(millis) * time.Millisecond, Enabled: false, Job: PrintTask}
	it.Start()
}

//export StopTimer
func StopTimer() {
	it.Stop()
}

func main() {}
