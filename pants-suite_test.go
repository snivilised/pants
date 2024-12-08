package pants_test

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/pants"
)

func TestPants(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pants Suite")
}

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
)

const (
	Param           = 100
	AntsSize        = 1000
	TestSize        = 10000
	n               = 100000
	PoolSize        = 10
	InputBufferSize = 3
)

const (
	RunTimes           = 1e6
	PoolCap            = 5e4
	BenchParam         = 10
	DefaultExpiredTime = 10 * time.Second
	CheckCloseInterval = time.Second / 100
	TimeoutOnSend      = time.Second
)

var curMem uint64

func demoFunc() {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}

func demoPoolFunc(inputCh pants.InputEnvelope) {
	n, _ := inputCh.Param().(int)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func demoPoolManifoldFunc(input int) (int, error) {
	time.Sleep(time.Duration(input) * time.Millisecond)

	return n + 1, nil
}

var stopLongRunningFunc int32

func longRunningFunc() {
	for atomic.LoadInt32(&stopLongRunningFunc) == 0 {
		runtime.Gosched()
	}
}

var stopLongRunningPoolFunc int32

func longRunningPoolFunc(arg pants.InputEnvelope) {
	if ch, ok := arg.Param().(chan struct{}); ok {
		<-ch
		return
	}
	for atomic.LoadInt32(&stopLongRunningPoolFunc) == 0 {
		runtime.Gosched()
	}
}

func ShowMemStats() {
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	GinkgoWriter.Printf("memory usage:%d MB", curMem)
}
