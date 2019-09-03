package main

/*
typedef struct _Test{
	int a; 
	char* b; 
} Test;
*/
import "C"

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

type Test struct {
	a int
	b string
}

var count int
var mtx sync.Mutex

//export Add
func Add(a, b int) int { return a + b }

//export Cosine
func Cosine(x float64) float64 { return math.Cos(x) }

//export Sort
func Sort(vals [] int) { sort.Ints(vals) }

//export CreateTest
func CreateTest(a int, b string) C.Test {
	return C.Test{
		a: C.int(a),
		b: C.CString(b),
	}
}

//export Log
func Log(msg string) int {
  mtx.Lock()
  defer mtx.Unlock()
  fmt.Println(msg)
  count++
  return count
}

// export Increment
func (t *Test) Increment() {
	t.a = t.a + 1
}

// export SetString
func (t *Test) SetString(newString string) (currentString string) {
	currentString = t.b
	t.b = newString
	return
}

func main() {}
