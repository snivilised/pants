// Package main example program
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/snivilised/pants"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var wg sync.WaitGroup

	// 1. Create a ShellPool that uses bash -i
	// In a real scenario, this would load the user's .bashrc once per worker.
	fmt.Println("🚀 Initializing ShellPool with 4 persistent bash workers...")
	pool, err := pants.NewShellPool(ctx, "bash", &wg,
		pants.WithSize(4),
		pants.WithOutput(100, time.Millisecond*100, time.Second),
	)
	if err != nil {
		panic(err)
	}

	// 2. Start observing results
	go func() {
		for output := range pool.Observe() {
			if output.Error != nil {
				fmt.Printf("❌ Job %s failed: %v\n", output.ID, output.Error)
			} else {
				fmt.Printf("✅ Job %s (Worker %v) Result: %s\n", 
					output.ID, output.WorkerID, output.Payload)
			}
		}
	}()

	// 3. Submit jobs
	// These jobs will run inside the persistent shell sessions.
	// Notice we are just passing the command string.
	commands := []string{
		"echo 'Hello from persistent shell!'",
		"uptime",
		"whoami",
		"date",
		"echo $SHELL",
		"export PANTS_VAR='Look, I can persist state!'",
		"echo $PANTS_VAR", // This will only work if it hits the same worker!
	}

	fmt.Println("📬 Submitting jobs...")
	for _, cmd := range commands {
		_ = pool.Post(ctx, cmd)
	}

	// 4. Conclude and wait
	fmt.Println("🏁 Concluding pool and waiting for results...")
	pool.Conclude(ctx)
	wg.Wait()
	
	pool.Release(ctx)
	fmt.Println("👋 Done.")
}
