package pprof

import (
	"errors"
	"os"
	"runtime/pprof"
	"sync/atomic"
)

var (
	cpuProfiling          int32
	cpuProfile            *os.File
	ErrCPUProfileStart    = errors.New("CPU profile already start")
	ErrCPUProfileNotStart = errors.New("CPU profile not start")
)

func StartCPUProfile(file string) error {
	if atomic.CompareAndSwapInt32(&cpuProfiling, 0, 1) {
		cpuProfile, err := os.Create(file)
		if err != nil {
			return err
		}
		return pprof.StartCPUProfile(cpuProfile)
	}
	return nil
}

func StopCPUProfile() {
	if atomic.LoadInt32(&cpuProfiling) == 1 {
		pprof.StopCPUProfile()
		cpuProfile.Close()
		cpuProfile = nil
	}
}

// Invoke runtime/pprof.Lookup(name, debug) then save to file.
//
//	goroutine    - stack traces of all current goroutines
//	heap         - a sampling of all heap allocations
//	threadcreate - stack traces that led to the creation of new OS threads
//	block        - stack traces that led to blocking on synchronization primitives
//
func SaveProfile(name, file string, debug int) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return pprof.Lookup(name).WriteTo(f, debug)
}
