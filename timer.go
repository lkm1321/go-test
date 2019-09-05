package main

/*
#cgo CFLAGS: -std=c99
#include <stdlib.h>
#include <stdbool.h> 

typedef struct _Timer{
	int millis; 
	char* message; 
	bool enabled; 
} Timer; 

typedef Timer* TimerPtr; 
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
	"sync"
)

/*
Below is a copy-and-paste from the https://github.com/mattn/go-pointer project. 
It's used to pass GoPointer to a Go object to the C space, so that the Go object persists throughout the C program. 
Doing so allows us to convert a C object to Go object by refernce, so the transparency is improved. 
*/
var (
	mutex sync.Mutex
	store = map[unsafe.Pointer]interface{}{}
)

func Save(v interface{}) unsafe.Pointer {
	if v == nil {
		return nil
	}

	// Generate real fake C pointer.
	// This pointer will not store any data, but will bi used for indexing purposes.
	// Since Go doest allow to cast dangling pointer to unsafe.Pointer, we do rally allocate one byte.
	// Why we need indexing, because Go doest allow C code to store pointers to Go data.
	var ptr unsafe.Pointer = C.malloc(C.size_t(1))
	if ptr == nil {
		panic("can't allocate 'cgo-pointer hack index pointer': ptr == nil")
	}

	mutex.Lock()
	store[ptr] = v
	mutex.Unlock()

	return ptr
}

func Restore(ptr unsafe.Pointer) (v interface{}) {
	if ptr == nil {
		return nil
	}

	mutex.Lock()
	v = store[ptr]
	mutex.Unlock()
	return
}

func Unref(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	mutex.Lock()
	delete(store, ptr)
	mutex.Unlock()

	C.free(ptr)
}

/*
Presented below is an approach for interfacing a C object to a Go object in a transparent manner. 
The GoObject is really just a reference to the orinal C object. 
To use this part, compile make ctimer, and run ./ctimer
*/
type Timer struct {
	_Timer C.TimerPtr
}

func (t *Timer) Cobj() C.Timer {
	return (*t._Timer)
}

func (t *Timer) Enabled() bool {
	return bool(t.Cobj().enabled)
}

func (t *Timer) GetString() string {
	return C.GoString(t.Cobj().message)
}

func (t *Timer) GetMillis() int {
	return int(t.Cobj().millis)
}

func (t *Timer) Isr() {
	if bool( t.Enabled() ) {
		fmt.Println(t.GetString())
		time.AfterFunc(time.Duration(t.GetMillis())*time.Millisecond, t.Isr)
	}
}

//export InitGoTimer
func InitGoTimer(CTimer C.TimerPtr) (unsafe.Pointer) {
	return Save(&Timer{_Timer: CTimer}) 
}

//export StartGoTimer
func StartGoTimer(goTimerUnsafe unsafe.Pointer) {
	var goTimer *Timer = Restore(goTimerUnsafe).(*Timer)
	time.AfterFunc(time.Duration(goTimer.Cobj().millis)*time.Millisecond, goTimer.Isr)
}

/*
Presented below is an approach where we copy a C object to Go object by value. 
The C object is only involved in the creation of Go object. 
The Go object must be a global in the Go space, so that it persists after conversion. 
Then, extern functions are used to perform operations on the Go object. 
To use this example. compile make timer, and run ./main
*/
type IntervalTimer struct {
	Interval time.Duration
	Enabled bool
	Message string
}

func (it *IntervalTimer) Isr() {
	if it.Enabled {
		fmt.Println(it.Message)
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

func (it *IntervalTimer) Toggle() {
	it.Enabled = !it.Enabled
	if (it.Enabled) {
		it.Start()
	}
}

// this needs to be global.
// had it been defined inside StartTimer, it would go out of scope after StartTimer call.
var it *IntervalTimer

// C.Timer gets defined by the comment block above import "C"
// this can be in seperate header file if needed. 

//export StartTimer
func StartTimer(timer C.Timer) {
	it = &IntervalTimer{Interval: time.Duration(timer.millis) * time.Millisecond, Enabled: false, Message: C.GoString(timer.message)}
	it.Start()
}

//export StopTimer
func StopTimer() {
	it.Stop()
}

//export ToggleTimer
func ToggleTimer() {
	fmt.Println("Toggling!")
	it.Toggle()
}

func main() {}
