package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/snivilised/pants"
	"github.com/snivilised/pants/internal/third/lo"
)

// Demonstrates Timeout On Send
func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("⏱️ timeout on send: '%.2f's\n", TimeoutOnSend.Seconds())
	pool, err := pants.NewManifoldFuncPool(
		ctx, func(input int) (int, error) {
			time.Sleep(time.Duration(input) * time.Millisecond)

			return n + 1, nil
		}, &wg,
		pants.WithSize(AntsSize),
		pants.WithOutput(OutputChSize, CheckCloseInterval, TimeoutOnSend),
	)

	defer pool.Release(ctx)
	if cc := pool.CancelCh(); cc != nil {
		pants.StartCancellationMonitor(ctx, cancel, &wg, cc, func() {
			fmt.Print("🔴 cancellation received, cancelling...\n")
		})
	}

	if err != nil {
		fmt.Printf("🔥 error creating pool: '%v'\n", err)
		return
	}

	wg.Add(1)
	go produce(ctx, pool, &wg) //nolint:wsl // pedantic

	wg.Add(1)
	go consume(ctx, pool, &wg) //nolint:wsl // pedantic

	fmt.Printf("pool with func, no of running workers:%d\n",
		pool.Running(),
	)
	wg.Wait()
	fmt.Println("🏁 (manifold-func-pool, timeout on send) FINISHED")
}

const (
	AntsSize           = 1000
	n                  = 100000
	OutputChSize       = 10
	Param              = 100
	CheckCloseInterval = time.Second / 10
	TimeoutOnSend      = time.Second
)

func produce(ctx context.Context,
	pool *pants.ManifoldFuncPool[int, int],
	wg pants.WaitGroup,
) {
	defer wg.Done()

	for i, n := 0, 100; i < n; i++ {
		_ = pool.Post(ctx, Param)
	}

	// required to inform the worker pool that no more jobs will be submitted.
	// failure to invoke Conclude will result in a never ending worker pool.
	//
	pool.Conclude(ctx)
}

func consume(ctx context.Context,
	pool *pants.ManifoldFuncPool[int, int],
	wg pants.WaitGroup,
) {
	defer wg.Done()

	const (
		fast    = time.Second / 10
		slow    = time.Second * 3
		barrier = 10
	)

	for count := 0; ; count++ {
		// Slow consumer after 10 iterations, resulting in a timeout
		//
		time.Sleep(
			lo.Ternary(count > barrier, slow, fast),
		)

		// NB: can not range over the observe channel since range
		// is non-preempt-able and therefore does not react to
		// ctx.Done.
		select {
		case output := <-pool.Observe():
			fmt.Printf("🍒 payload: '%v', id: '%v', seq: '%v' (e: '%v')\n",
				output.Payload, output.ID, output.SequenceNo, output.Error,
			)
		case <-ctx.Done():
			return
		}
	}
}
