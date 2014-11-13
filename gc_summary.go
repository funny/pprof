package care

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
	"time"
)

var startTime = time.Now()

func GCSummary(writer io.Writer) {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	gcstats := &debug.GCStats{PauseQuantiles: make([]time.Duration, 100)}
	debug.ReadGCStats(gcstats)

	if gcstats.NumGC > 0 {
		lastPause := gcstats.Pause[0]
		elapsed := time.Now().Sub(startTime)
		overhead := float64(gcstats.PauseTotal) / float64(elapsed) * 100
		allocatedRate := float64(memStats.TotalAlloc) / elapsed.Seconds()

		fmt.Fprintf(writer,
			"NumGC: %d, LastPause: %v, Pause(Avg): %v, Overhead: %3.2f%%, Alloc: %s, Sys: %s, Alloc(Rate): %s/s, Histogram: %v %v %v\n",
			gcstats.NumGC,
			lastPause,
			durationAvg(gcstats.Pause),
			overhead,
			readableSize(memStats.Alloc),
			readableSize(memStats.Sys),
			readableSize(uint64(allocatedRate)),
			gcstats.PauseQuantiles[94],
			gcstats.PauseQuantiles[98],
			gcstats.PauseQuantiles[99],
		)
	} else {
		elapsed := time.Now().Sub(startTime)
		allocatedRate := float64(memStats.TotalAlloc) / elapsed.Seconds()

		fmt.Fprintf(writer,
			"NumGC: 0, Alloc: %s, Sys:%s, Alloc(Rate): %s/s\n",
			readableSize(memStats.Alloc),
			readableSize(memStats.Sys),
			readableSize(uint64(allocatedRate)),
		)
	}
}

func durationAvg(items []time.Duration) time.Duration {
	var sum time.Duration
	for _, item := range items {
		sum += item
	}
	return time.Duration(int64(sum) / int64(len(items)))
}

func readableSize(bytes uint64) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%d B", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.2f K", float64(bytes)/1024)
	case bytes < 1024*1024*1024:
		return fmt.Sprintf("%.2f M", float64(bytes)/1024/1024)
	default:
		return fmt.Sprintf("%.2f G", float64(bytes)/1024/1024/1024)
	}
}
