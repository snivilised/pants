package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/snivilised/pants"
)

// Demonstrates use of manifold func base worker pool where
// the client manifold func returns an output and an error.
// Submission to the pool occurs via an input channel as opposed
// directly invoking Post on the pool.

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pants.NewManifoldFuncPool(
		ctx, func(input int) (int, error) {
			time.Sleep(time.Duration(input) * time.Millisecond)

			return n + 1, nil
		}, &wg,
		pants.WithSize(AntsSize),
		pants.WithInput(InputChSize),
		pants.WithOutput(OutputChSize, CheckCloseInterval, TimeoutOnSend),
	)

	defer pool.Release(ctx)

	if err != nil {
		fmt.Printf("🔥 error creating pool: '%v'\n", err)
		return
	}

	wg.Add(1)

	go inject(ctx, pool, &wg)

	wg.Add(1)

	go consume(ctx, pool, &wg)

	fmt.Printf("pool with func, no of running workers:%d\n",
		pool.Running(),
	)
	wg.Wait()
	fmt.Println("🏁 (manifold-func-pool, input injected via channel) FINISHED")
}

const (
	AntsSize           = 1000
	n                  = 100000
	InputChSize        = 10
	OutputChSize       = 10
	Param              = 100
	OutputChTimeout    = time.Second / 2 // do not use a value that is similar to interval
	CheckCloseInterval = time.Second / 10
	TimeoutOnSend      = time.Second * 2
)

func inject(ctx context.Context,
	pool *pants.ManifoldFuncPool[int, int],
	wg pants.WaitGroup,
) {
	defer wg.Done()

	inputCh := pool.Source(ctx, wg)
	for i, n := 0, 100; i < n; i++ {
		inputCh <- Param
	}

	// required to inform the worker pool that no more jobs will be submitted.
	// failure to close the input channel will result in a never ending
	// worker pool.
	//
	close(inputCh)
}

func consume(_ context.Context,
	pool *pants.ManifoldFuncPool[int, int],
	wg pants.WaitGroup,
) {
	defer wg.Done()

	// We don't need to use a timeout on the observe channel
	// because our producer invokes Conclude, which results in
	// the observe channel being closed, terminating the range.
	// This aspect is specific to this example and clients may
	// have to use different strategies depending on their use-case,
	// eg support for context cancellation.
	//
	for output := range pool.Observe() {
		fmt.Printf("🍒 payload: '%v', output-id: '%v', seq: '%v' worker-id: '%v', (e: '%v')\n",
			output.Payload, output.ID, output.SequenceNo, output.WorkerID, output.Error,
		)
	}
}
