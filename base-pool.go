package pants

import (
	"sync/atomic"

	"github.com/snivilised/pants/locale"
)

type (
	basePool[I, O any] struct {
		wg         WaitGroup
		sequence   int32
		inputDupCh *Duplex[I]
		oi         *outputInfo[O]
		ending     bool
	}
)

func (p *basePool[I, O]) next() int32 {
	return atomic.AddInt32(&p.sequence, int32(1))
}

// Observe returns a channel which can be read from to obtain
// the output of the pool. Using Observe here is only ever valid
// if an output has been requested using the WithOutput operator.
// If WithOutput has not been called but the client invokes Observe,
// then this is characterised as a serious error and a panic occurs
// A panic is required in this situation to allow the client to
// range over the returned channel, save in the knowledge that it is
// indeed valid.
func (p *basePool[I, O]) Observe() JobOutputStreamR[O] {
	if p.oi == nil {
		panic(locale.ErrBadObservation)
	}

	return p.oi.outputDupCh.ReaderCh
}

// CancelCh
func (p *basePool[I, O]) CancelCh() CancelStreamR {
	if p.oi != nil {
		return p.oi.cancelDupCh.ReaderCh
	}

	return nil
}
