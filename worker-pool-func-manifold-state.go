package pants

import (
	"context"

	"github.com/snivilised/pants/internal/third/ants"
)

type (
	// ManifoldStateFunc is the pre-defined function registered with the worker
	// pool, executed for each incoming job with worker-local state.
	ManifoldStateFunc[I, O, S any] func(input I, state S) (O, error)
)

// ManifoldStatePool is a wrapper around the underlying ants function based
// worker pool with support for per-worker state.
type ManifoldStatePool[I, O, S any] struct {
	basePool[I, O]
	functionalPool
}

// NewManifoldStatePool creates a new manifold state based worker pool.
func NewManifoldStatePool[I, O, S any](ctx context.Context,
	mf ManifoldStateFunc[I, O, S],
	wg WaitGroup,
	options ...Option,
) (*ManifoldStatePool[I, O, S], error) {
	var (
		oi *outputInfo[O]
		wi *outputInfoW[O]
		o  = ants.NewOptions(options...)
	)

	if oi = newOutputInfo[O](o); oi != nil {
		wi = fromOutputInfo(o, oi)
	}

	pool, err := ants.NewPoolWithFunc(ctx, func(input InputEnvelope) {
		manifoldStateFuncResponse(ctx, mf, input, wi)
	}, ants.WithOptions(*o))

	return &ManifoldStatePool[I, O, S]{
		basePool: basePool[I, O]{
			wg: wg,
			oi: oi,
		},
		functionalPool: functionalPool{
			pool: pool,
		},
	}, err
}

// Post allows the client to submit to the work pool represented by
// input values of type I.
func (p *ManifoldStatePool[I, O, S]) Post(ctx context.Context, input I) error {
	o := p.pool.GetOptions()
	job := Job[I]{
		ID:         o.Generator.Generate(),
		Input:      input,
		SequenceNo: int(p.next()),
	}

	return p.pool.Invoke(ctx, job)
}

// Source returns an input stream through which the client can submit
// jobs to the pool.
func (p *ManifoldStatePool[I, O, S]) Source(ctx context.Context,
	wg WaitGroup,
) SourceStreamW[I] {
	o := p.pool.GetOptions()

	p.inputDupCh = source(ctx, wg, o,
		injector[I](func(input I) error {
			return p.Post(ctx, input)
		}),
		terminator(func() {
			p.Conclude(ctx)
		}),
	)

	return p.inputDupCh.WriterCh
}

// Conclude signifies to the worker pool that no more work will be submitted.
func (p *ManifoldStatePool[I, O, S]) Conclude(ctx context.Context) {
	conclude[I, O](ctx, &p.basePool, &p.functionalPool)
}

func manifoldStateFuncResponse[I, O, S any](ctx context.Context,
	mf ManifoldStateFunc[I, O, S],
	input InputEnvelope,
	wi *outputInfoW[O],
) {
	if job, ok := input.Param().(Job[I]); ok {
		var state S
		if s, ok := input.State().(S); ok {
			state = s
		}

		payload, e := mf(job.Input, state)

		if wi != nil {
			_ = respond(ctx, wi, &JobOutput[O]{
				ID:         job.ID,
				SequenceNo: job.SequenceNo,
				Payload:    payload,
				Error:      e,
				WorkerID:   input.WorkerID(),
			})
		}
	}
}
