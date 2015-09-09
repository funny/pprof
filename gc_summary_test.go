package pprof

import (
	"runtime"
	"testing"
)

func Test_GCSummary(t *testing.T) {
	summary := GCSummary()
	println(" ", summary.String())

	runtime.GC()

	summary = GCSummary()
	println(" ", summary.String())

	println(" ", GCSummaryColumns)
	println(" ", summary.CSV())
}
