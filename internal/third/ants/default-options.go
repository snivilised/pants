package ants

import (
	"runtime"
	"time"
)

var (
	// DefaultOutputOptions denotes the default output options used if the client
	// requests an output channel without having defined an output in the options.
	// This is not the recommended way using the output as the defaults here are
	// not likely to be suitable for their use-case, but we do this in order not to
	// panic or forced to return an error.
	DefaultOutputOptions *OutputOptions
)

const (
	scale = 5
)

func init() {
	DefaultOutputOptions = &OutputOptions{
		BufferSize:         uint(runtime.NumCPU()),
		CheckCloseInterval: time.Second * scale,
		TimeoutOnSend:      time.Second * scale,
	}
}
