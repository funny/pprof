package pprof

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

var startTime = time.Now()

// Shorthand for GCSummary().Save(file)
func SaveGCSummary(file string) error {
	return GCSummary().Save(file)
}

// Shorthand for GCSummary().SaveCSV(file)
func SaveGCSummaryGCV(file string) error {
	return GCSummary().SaveCSV(file)
}

type GCSummaryInfo struct {
	NumGC      int64
	LastPause  time.Duration
	PauseAvg   time.Duration
	Overhead   float64
	Alloc      uint64
	Sys        uint64
	AllocRate  uint64
	Histogram1 time.Duration
	Histogram2 time.Duration
	Histogram3 time.Duration
}

// Get GC summary.
func GCSummary() *GCSummaryInfo {
	gcstats := debug.GCStats{PauseQuantiles: make([]time.Duration, 100)}
	debug.ReadGCStats(&gcstats)

	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	elapsed := time.Now().Sub(startTime)

	summary := &GCSummaryInfo{
		Alloc:     memStats.Alloc,
		Sys:       memStats.Sys,
		AllocRate: uint64(float64(memStats.TotalAlloc) / elapsed.Seconds()),
	}

	if gcstats.NumGC > 0 {
		summary.NumGC = gcstats.NumGC
		summary.LastPause = gcstats.Pause[0]
		summary.PauseAvg = durationAvg(gcstats.Pause)
		summary.Overhead = float64(gcstats.PauseTotal) / float64(elapsed) * 100
		summary.Histogram1 = gcstats.PauseQuantiles[94]
		summary.Histogram2 = gcstats.PauseQuantiles[98]
		summary.Histogram3 = gcstats.PauseQuantiles[99]
	}

	return summary
}

// Humman readable string.
func (summary *GCSummaryInfo) String() string {
	buffer := new(bytes.Buffer)
	summary.Write(buffer)
	return buffer.String()
}

// CSV string.
func (summary *GCSummaryInfo) CSV() string {
	buffer := new(bytes.Buffer)
	summary.WriteCSV(buffer)
	return buffer.String()
}

// Write as humman readable format.
func (summary *GCSummaryInfo) Write(writer io.Writer) error {
	_, err := fmt.Fprintf(writer,
		"NumGC: %d, LastPause: %v, Pause(Avg): %v, Overhead: %3.2f%%, Alloc: %s, Sys: %s, Alloc(Rate): %s/s, Histogram: %v %v %v\n",
		summary.NumGC,
		summary.LastPause,
		summary.PauseAvg,
		summary.Overhead,
		formatSize(summary.Alloc),
		formatSize(summary.Sys),
		formatSize(summary.AllocRate),
		summary.Histogram1,
		summary.Histogram2,
		summary.Histogram3,
	)
	return err
}

// Save as human readable file.
func (summary *GCSummaryInfo) Save(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return summary.Write(f)
}

// GC summary CSV column names.
const GCSummaryColumns = "NumGC,LastPause,Pause(Avg),Overhead,Alloc,Sys,Alloc(Rate),Histogram1,Histogram2,Histogram3"

// Write as CSV format.
func (summary *GCSummaryInfo) WriteCSV(writer io.Writer) error {
	_, err := fmt.Fprintf(writer,
		"%d,%d,%d,%3.2f,%d,%d,%d,%d,%d,%d\n",
		summary.NumGC,
		summary.LastPause,
		summary.PauseAvg,
		summary.Overhead,
		summary.Alloc,
		summary.Sys,
		summary.AllocRate,
		summary.Histogram1,
		summary.Histogram2,
		summary.Histogram3,
	)
	return err
}

// Save as CSV file.
func (summary *GCSummaryInfo) SaveCSV(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return summary.WriteCSV(f)
}

func durationAvg(items []time.Duration) time.Duration {
	var sum time.Duration
	for _, item := range items {
		sum += item
	}
	return time.Duration(int64(sum) / int64(len(items)))
}

func formatSize(bytes uint64) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%dB", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.2fK", float64(bytes)/1024)
	case bytes < 1024*1024*1024:
		return fmt.Sprintf("%.2fM", float64(bytes)/1024/1024)
	default:
		return fmt.Sprintf("%.2fG", float64(bytes)/1024/1024/1024)
	}
}
