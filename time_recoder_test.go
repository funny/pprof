package overall

import (
	"bytes"
	"testing"
	"time"
)

func Test_TimeRecoder(t *testing.T) {
	recoder := NewTimeRecoder()

	for i := 0; i < 10; i++ {
		t := time.Now()
		do_something(100)
		recoder.Record("do_something(100)", time.Since(t))
	}

	for i := 0; i < 10; i++ {
		t := time.Now()
		do_something(1000)
		recoder.Record("do_something(1000)", time.Since(t))
	}

	for i := 0; i < 10; i++ {
		t := time.Now()
		do_something(10000)
		recoder.Record("do_something(10000)", time.Since(t))
	}

	for i := 0; i < 10; i++ {
		t := time.Now()
		do_something(100000)
		recoder.Record("do_something(100000)", time.Since(t))
	}

	buffer := new(bytes.Buffer)
	recoder.WriteCSV(buffer)
	println()
	println(buffer.String())
	println()
}

func do_something(n int) int {
	m := 0
	for i := 0; i < n; i++ {
		m += i
	}
	return m
}
