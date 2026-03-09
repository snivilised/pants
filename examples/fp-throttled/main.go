// Package main example program
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/snivilised/pants"
)

// Demonstrates that when all workers are engaged and the pool is at capacity,
// new incoming jobs are blocked, until a worker becomes free. The invoked function
// takes a second to complete. The PRE and POST indicators reflect this:
//
// PRE: <--- (n: 0) [13:56:22] 🍋
// => running: '0')
// POST: <--- (n: 0) [13:56:22] 🍊
// PRE: <--- (n: 1) [13:56:22] 🍋
// => running: '1')
// POST: <--- (n: 1) [13:56:22] 🍊
// PRE: <--- (n: 2) [13:56:22] 🍋
// => running: '2')
// POST: <--- (n: 2) [13:56:22] 🍊
// PRE: <--- (n: 3) [13:56:22] 🍋
// => running: '3')
// <--- (n: 2)🍒
// <--- (n: 1)🍒
// <--- (n: 0)🍒
// <--- (n: 3)🍒
// POST: <--- (n: 3) [13:56:23] 🍊
//
// Considering the above, whilst the pool is not at capacity, each new submission is
// executed immediately, as a new worker can be allocated to those jobs (n=0..2).
// Once the pool has reached capacity (n=3), the PRE is blocked, because its corresponding
// POST doesn't happen until a second later; this illustrates the blocking.
//

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const NoW = 3

	pool, _ := pants.NewFuncPool[int, int](ctx, func(input pants.InputEnvelope) {
		n, _ := input.Param().(int)
		id := input.WorkerID()
		fmt.Printf("<--- (n: %v, from: %v)🍒 \n", n, id)
		time.Sleep(time.Second)
	}, &wg,
		pants.WithSize(NoW),
		pants.WithNonblocking(false),
	)

	defer pool.Release(ctx)

	for i := 0; i < 30; i++ { // producer
		fmt.Printf("PRE: <--- (n: %v) [%v] 🍋 \n=> running: '%v'\n",
			i, time.Now().Format(time.TimeOnly), pool.Running(),
		)
		_ = pool.Post(ctx, i)
		fmt.Printf("POST: <--- (n: %v) [%v] 🍊 \n", i, time.Now().Format(time.TimeOnly))
	}

	fmt.Printf("pool with func, no of running workers:%d\n",
		pool.Running(),
	)

	// Note, we don't need to inform the pool of the end of the workload
	// since this pool is not emitting output.
	wg.Wait()
	fmt.Println("🏁 (func-pool) FINISHED")
}
