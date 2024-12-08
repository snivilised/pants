package pants

const (
	MaxWorkers = 100
)

type (
	// SourceStream bi-directional channel of item I. The source
	// stream represents the input stream through which the jobs
	// can be submitted to the worker pool.
	SourceStream[I any] chan I

	// SourceStreamR is the read side of the SourceStream
	SourceStreamR[I any] <-chan I

	// SourceStreamW is the write side of the SourceStream
	SourceStreamW[I any] chan<- I

	// Job of input item I that can is submitted to the worker pool
	Job[I any] struct {
		// ID uniquely identifies the Job
		ID string

		// SequenceNo represents the order of the Job
		SequenceNo int

		// Input source item of the Job
		Input I
	}

	// JobOutput represents the output of Job execution
	JobOutput[O any] struct {
		ID         string
		SequenceNo int
		Payload    O
		Error      error
		WorkerID   RoutineID
	}

	// JobStream bi-directional channel of Jobs of I
	JobStream[I any] chan Job[I]

	// JobStreamR is the read side of the JobStream
	JobStreamR[I any] <-chan Job[I]

	// JobStreamW is the write side of the JobStream
	JobStreamW[I any] chan<- Job[I]

	// JobOutputStream bi-directional channel of JobOutput of O
	JobOutputStream[O any] chan JobOutput[O]

	// JobOutputStreamR is the read side of the JobOutputStream
	JobOutputStreamR[O any] <-chan JobOutput[O]

	// JobOutputStreamW is the write side of the JobOutputStream
	JobOutputStreamW[O any] chan<- JobOutput[O]

	// Duplex represents a channel with multiple views, to be used
	// by clients that need to hand out different ends of the same
	// channel to different entities.
	Duplex[T any] struct {
		Channel  chan T
		ReaderCh <-chan T
		WriterCh chan<- T
	}

	// DuplexJobOutput defines an Duplex of JobOutput of O
	DuplexJobOutput[O any] Duplex[JobOutput[O]]

	// CancelWorkSignal item send to cancel indication
	CancelWorkSignal struct{}

	// CancelStream bi-directional channel of CancelWorkSignal
	CancelStream = chan CancelWorkSignal

	// CancelStreamR is the read side of the CancelStream
	CancelStreamR = <-chan CancelWorkSignal

	// CancelStreamW is the write side of the CancelStream
	CancelStreamW = chan<- CancelWorkSignal

	// OnCancel is the callback required by StartCancellationMonitor
	OnCancel func()

	// WaitGroup allows the core sync.WaitGroup to be decorated by the client
	// for debugging purposes.
	WaitGroup interface {
		Add(delta int)
		Done()
		Wait()
	}
)

// ExecutiveFunc is the function executed by the worker pool for each
// submitted job. Each job is characterised by its input I and its
// output O.
type ExecutiveFunc[I, O any] func(j Job[I]) (JobOutput[O], error)

// Invoke use by the worker pool to execute the job
func (f ExecutiveFunc[I, O]) Invoke(j Job[I]) (JobOutput[O], error) {
	return f(j)
}

// NewDuplex creates a new instance of a Duplex with all members populated
func NewDuplex[T any](channel chan T) *Duplex[T] {
	return &Duplex[T]{
		Channel:  channel,
		ReaderCh: channel,
		WriterCh: channel,
	}
}
