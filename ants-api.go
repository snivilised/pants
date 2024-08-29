package pants

import "github.com/snivilised/pants/internal/third/ants"

type (
	// IDGenerator is a sequential unique id generator interface
	IDGenerator = ants.IDGenerator

	// InputParam
	InputParam = ants.InputParam

	// Option represents the ants functional option.
	Option = ants.Option

	// Option represents the ants options.
	Options = ants.Options

	// PoolFunc ants pool function
	PoolFunc = ants.PoolFunc

	// Sequential represents te ants sequential ID generator
	Sequential = ants.Sequential

	// TaskFunc ants task function
	TaskFunc = ants.TaskFunc
)

var (
	// WithDisablePurge indicates whether we turn off automatically purge
	WithDisablePurge = ants.WithDisablePurge

	// WithExpiryDuration sets up the interval time of cleaning up goroutines
	WithExpiryDuration = ants.WithExpiryDuration

	// WithGenerator sets up an ID generator
	WithGenerator = ants.WithGenerator

	// WithInput sets input buffer size
	WithInput = ants.WithInput

	// WithMaxBlockingTasks sets up the maximum number of goroutines that are
	// blocked when it reaches the capacity of pool.
	WithMaxBlockingTasks = ants.WithMaxBlockingTasks

	// WithNonblocking indicates that pool will return nil when there is no
	// available workers.
	WithNonblocking = ants.WithNonblocking

	// WithOptions accepts the whole options config.
	WithOptions = ants.WithOptions

	// WithOutput sets output characteristics:
	// size uint: defines the size of the output channel
	// interval time.Duration: usee by Conclude to check if its safe to close
	// the output channel, periodically, which is implemented within another Go routine.
	// timeout time.Duration: denotes the timeout used when the pool attempts
	// to send to the output channel
	WithOutput = ants.WithOutput

	// WithPanicHandler sets up panic handler.
	WithPanicHandler = ants.WithPanicHandler

	// WithPreAlloc indicates whether it should malloc for workers.
	WithPreAlloc = ants.WithPreAlloc

	// WithSize denotes the number of workers in the pool. Defaults
	// to number of CPUs available.
	WithSize = ants.WithSize
)
