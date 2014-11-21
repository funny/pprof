package overall

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"
)

type TimeRecoder struct {
	mutex   sync.RWMutex
	records map[string]*timeRecord
}

type timeRecord struct {
	Times         int64
	TotalUsedTime time.Duration
	MaxUsedTime   time.Duration
	MinUsedTime   time.Duration
}

func NewTimeRecoder() *TimeRecoder {
	return &TimeRecoder{
		records: make(map[string]*timeRecord),
	}
}

func (tr *TimeRecoder) Record(name string, usedTime time.Duration) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	r, exists := tr.records[name]

	if exists {
		r.Times += 1
		r.TotalUsedTime += usedTime

		if r.MaxUsedTime < usedTime {
			r.MaxUsedTime = usedTime
		}

		if r.MinUsedTime > usedTime {
			r.MinUsedTime = usedTime
		}
	} else {
		r = &timeRecord{1, usedTime, usedTime, usedTime}
		tr.records[name] = r
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
