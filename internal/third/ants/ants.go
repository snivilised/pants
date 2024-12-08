// MIT License

// Copyright (c) 2018 Andy Pan

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package ants

import (
	"log"
	"math"
	"os"
	"runtime"
	"time"
)

const (
	// DefaultAntsPoolSize is the default capacity for a default goroutine pool.
	DefaultAntsPoolSize = math.MaxInt32

	// DefaultCleanIntervalTime is the interval time to clean up goroutines.
	DefaultCleanIntervalTime = time.Second

	releaseTimeoutInterval = 10
)

const (
	// OPENED represents that the pool is opened.
	OPENED = iota

	// CLOSED represents that the pool is closed.
	CLOSED
)

var (
	// workerChanCap determines whether the channel of a worker should be a buffered channel
	// to get the best performance. Inspired by fasthttp at
	// https://github.com/valyala/fasthttp/blob/master/workerpool.go#L139
	workerChanCap = func() int {
		// Use blocking channel if GOMAXPROCS=1.
		// This switches context from sender to receiver immediately,
		// which results in higher performance (under go1.5 at least).
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}

		// Use non-blocking workerChan if GOMAXPROCS>1,
		// since otherwise the sender might be dragged down if the receiver is CPU-bound.
		return 1
	}()

	// log.Lmsgprefix is not available in go1.13, just make an identical value for it.
	logLmsgprefix = 64
	defaultLogger = Logger(log.New(os.Stderr, "[ants]: ",
		log.LstdFlags|logLmsgprefix|log.Lmicroseconds),
	)
)

const nowTimeUpdateInterval = 500 * time.Millisecond

// Logger is used for logging formatted messages.
type Logger interface {
	// Printf must have the same semantics as log.Printf.
	Printf(format string, args ...interface{})
}

type (
	// RoutineID the identifier representing the underlying worker.
	RoutineID int32

	// WorkEnvelope the task wrapper that provides access to the
	// worker id allocated to each job.
	WorkEnvelope interface {
		WorkerID() RoutineID
	}

	// TaskFunc represents the job function executed by task based
	// worker pools.
	TaskFunc func()

	// TaskStream the channel of tasks processed by task based worker
	// pools.
	TaskStream chan *TaskEnvelope

	// InputParam the input passed to the function for func based
	// worker pools.
	InputParam interface{}

	// InputEnvelope the input wrapper with an input
	InputEnvelope interface {
		WorkEnvelope
		Param() InputParam
	}

	// PoolFunc represents the job function executed by func based
	// worker pools.
	PoolFunc func(InputEnvelope)

	// InputStream
	InputStream chan InputEnvelope

	// Nothing
	Nothing struct{}

	// Envelope is the underlying wrapper used for func based (with input)
	// worker pools.
	Envelope struct {
		ID    RoutineID
		Input interface{}
	}

	// TaskEnvelope is the underlying wrapper used for task based
	// worker pools.
	TaskEnvelope struct {
		ID   RoutineID
		Task TaskFunc
	}
)

func (e Envelope) WorkerID() RoutineID {
	return e.ID
}

func (e Envelope) Param() InputParam {
	return e.Input
}

func (e TaskEnvelope) WorkerID() RoutineID {
	return e.ID
}

func (e TaskEnvelope) Func() TaskFunc {
	return e.Task
}
