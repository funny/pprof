package tprof

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type P struct {
	mutex   sync.RWMutex
	records map[string]*record
}

type record struct {
	Times         int64
	TotalUsedTime time.Duration
	MaxUsedTime   time.Duration
	MinUsedTime   time.Duration
}

func New() *P {
	return &P{
		records: make(map[string]*record),
	}
}

func (p *P) Record(name string, usedTime time.Duration) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var r, exists = p.records[name]

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
		r = &record{1, usedTime, usedTime, usedTime}
		p.records[name] = r
	}
}

func (p *P) SaveFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	p.Save(file)
	return nil
}

func (p *P) Save(writer io.Writer) {
	results := p.getRecords()
	sort.Sort(results)

	buf := bufio.NewWriter(writer)

	buf.WriteString("name,times,avg,min,max,total\n")

	for _, r := range results {
		buf.WriteString(r.Name)
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.Times, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.AvgUsedTime, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.MinUsedTime, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.MaxUsedTime, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.TotalUsedTime, 10))
		buf.WriteByte('\n')
	}

	buf.Flush()
}

func (p *P) getRecords() resultList {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	results := make(resultList, 0, len(p.records))

	for name, d := range p.records {
		results = append(results, &result{
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

type result struct {
	Name          string
	Times         int64
	AvgUsedTime   int64
	MinUsedTime   int64
	MaxUsedTime   int64
	TotalUsedTime int64
}

type resultList []*result

func (this resultList) Len() int {
	return len(this)
}

func (this resultList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this resultList) Less(i, j int) bool {
	return this[i].AvgUsedTime > this[j].AvgUsedTime || (this[i].AvgUsedTime == this[j].AvgUsedTime && this[i].Times < this[j].Times)
}
