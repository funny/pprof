package overall

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type TimeRecoder struct {
	mutex   sync.RWMutex
	records map[string]*timeRecord
}

type timeRecord struct {
	Times int64 // total count

	// The used time represents the elapsed time as an int64 nanosecond count.
	TotalUsedTime int64
	MaxUsedTime   int64
	MinUsedTime   int64
}

func NewTimeRecoder() *TimeRecoder {
	return &TimeRecoder{
		records: make(map[string]*timeRecord),
	}
}

func (tr *TimeRecoder) Record(name string, usedTime time.Duration) {
	usedTimeNano := int64(usedTime)
	var r *timeRecord
	{
		tr.mutex.Lock()
		defer tr.mutex.Unlock()
		var exists bool
		r, exists = tr.records[name]
		if !exists {
			r = &timeRecord{1, usedTimeNano, usedTimeNano, usedTimeNano}
			tr.records[name] = r
			return
		}
	}

	atomic.AddInt64(&r.Times, 1)
	atomic.AddInt64(&r.TotalUsedTime, usedTimeNano)

	if r.MaxUsedTime < usedTimeNano {
		atomic.StoreInt64(&r.MaxUsedTime, usedTimeNano)
	}

	if r.MinUsedTime > usedTimeNano {
		atomic.StoreInt64(&r.MinUsedTime, usedTimeNano)
	}
}

func (tr *TimeRecoder) SaveCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return tr.WriteCSV(file)
}

func (tr *TimeRecoder) WriteCSV(writer io.Writer) error {
	results := tr.getRecords()
	sort.Sort(results)

	buf := bufio.NewWriter(writer)

	if _, err := fmt.Fprintln(writer, "name,times,avg,min,max,total"); err != nil {
		return err
	}

	for _, r := range results {
		if _, err := fmt.Fprintf(writer,
			"%s,%d,%d,%d,%d,%d\n",
			r.Name,
			r.Times,
			r.AvgUsedTime,
			r.MinUsedTime,
			r.MaxUsedTime,
			r.TotalUsedTime,
		); err != nil {
			return err
		}
	}

	return buf.Flush()
}

func (tr *TimeRecoder) getRecords() sortTimeRecords {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	results := make(sortTimeRecords, 0, len(tr.records))

	for name, d := range tr.records {
		results = append(results, &sortTimeRecord{
			Name:          name,
			Times:         d.Times,
			AvgUsedTime:   int64(d.TotalUsedTime) / d.Times,
			MaxUsedTime:   int64(d.MaxUsedTime),
			MinUsedTime:   int64(d.MinUsedTime),
			TotalUsedTime: int64(d.TotalUsedTime),
		})
	}

	return results
}

type sortTimeRecord struct {
	Name          string
	Times         int64
	AvgUsedTime   int64
	MinUsedTime   int64
	MaxUsedTime   int64
	TotalUsedTime int64
}

type sortTimeRecords []*sortTimeRecord

func (this sortTimeRecords) Len() int {
	return len(this)
}

func (this sortTimeRecords) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this sortTimeRecords) Less(i, j int) bool {
	return this[i].AvgUsedTime > this[j].AvgUsedTime || (this[i].AvgUsedTime == this[j].AvgUsedTime && this[i].Times < this[j].Times)
}
