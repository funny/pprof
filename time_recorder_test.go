package pprof

import (
	"bytes"
	"strconv"
	"testing"
	"time"
)

func Test_TimeRecorder(t *testing.T) {
	recoder := NewTimeRecorder()

	for i := 0; i < 10; i++ {
		t := time.Now()
		doSomething(100)
		recoder.Record("doSomething(100)", time.Since(t))
	}

	for i := 0; i < 10; i++ {
		t := time.Now()
		doSomething(1000)
		recoder.Record("doSomething(1000)", time.Since(t))
	}

	for i := 0; i < 10; i++ {
		t := time.Now()
		doSomething(10000)
		recoder.Record("doSomething(10000)", time.Since(t))
	}

	for i := 0; i < 10; i++ {
		t := time.Now()
		doSomething(100000)
		recoder.Record("doSomething(100000)", time.Since(t))
	}

	buffer := new(bytes.Buffer)
	recoder.WriteCSV(buffer)
	println()
	println(buffer.String())
	println()
}

func doSomething(n int) int {
	m := 0
	for i := 0; i < n; i++ {
		m += i
		for j := 0; j < 30; j++ {
			strconv.Itoa(m) // cost some CPU time
		}
	}
	return m
}
