// Utility functions and types in support of transit_tools, but not specific
// to that application. Much of util is focused on rate limited http fetching,
// and I/O of .tar and .csv files (compressed or not), including via goroutines
// to enable concurrent reading, decompressing, parsing and processing (and
// vice versa). With luck I'll come up with a partitioning into smaller
// packages.

package goioutil

import (
	"flag"
	"fmt"
	"runtime"
)

var (
	goMaxProcsFlag = flag.Int(
		"go-max-procs", runtime.NumCPU(),
		"Number of concurrent go routines to run concurrently")
)

func InitGOMAXPROCS() {
	cpus := runtime.NumCPU()
	fmt.Println("NumCPU:", cpus)
	max := runtime.GOMAXPROCS(*goMaxProcsFlag)
	fmt.Println("Original GOMAXPROCS:", max)
	max = runtime.GOMAXPROCS(-1)
	fmt.Println("Current GOMAXPROCS:", max)
}

func LogToStderrValue() (logtostderr, ok bool) {
	logtostderr_flag := flag.Lookup("logtostderr")
	if logtostderr_flag == nil {
		ok = false
		return
	}
	getter, ok := logtostderr_flag.Value.(flag.Getter)
	if !ok {
		return
	}
	logtostderr, ok = getter.Get().(bool)
	return
}

type CountingBitBucketWriter uint64

func (p *CountingBitBucketWriter) Write(data []byte) (int, error) {
	size := len(data)
	*p = *p + CountingBitBucketWriter(size)
	return size, nil
}

func (p *CountingBitBucketWriter) Size() uint64 {
	return uint64(*p)
}
