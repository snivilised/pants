package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/snivilised/pants"
)

// Demonstrates Timeout On Send as a result of not consuming the output
func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("⏱️ timeout on send, missing consumer: '%.2f's\n", TimeoutOnSend.Seconds())
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
	go produce(ctx, pool, &wg)

	fmt.Printf("pool with func, no of running workers:%d\n",
		pool.Running(),
	)
	wg.Wait()
	fmt.Println("🏁 (manifold-func-pool, missing consumer) FINISHED")
}

const (
	AntsSize           = 1000
	n                  = 100000
	OutputChSize       = 10
	Param              = 100
	CheckCloseInterval = time.Second / 10
	TimeoutOnSend      = time.Second * 3
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
