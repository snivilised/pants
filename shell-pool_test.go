package pants_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/pants"
)

type mockShellSession struct {
	closed bool
}

func (s *mockShellSession) Execute(ctx context.Context, command string) (string, error) {
	if s.closed {
		return "", fmt.Errorf("session closed")
	}
	return fmt.Sprintf("executed: %s", command), nil
}

func (s *mockShellSession) Close() error {
	s.closed = true
	return nil
}

var _ = Describe("ShellPool", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		wg     sync.WaitGroup
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
		wg.Wait()
	})

	It("should execute commands using pooled shell sessions", func() {
		poolSize := 2
		jobCount := 5

		// We use the specialized NewShellPool but we override the initializer
		// for testing to avoid actual shell processes.
		// Actually, let's just test the ShellPool with a mock session via a helper.

		mf := func(command string, session pants.ShellSession) (string, error) {
			return session.Execute(ctx, command)
		}

		initializer := func(id pants.RoutineID) interface{} {
			return &mockShellSession{}
		}

		finalizer := func(state interface{}) {
			if s, ok := state.(pants.ShellSession); ok {
				_ = s.Close()
			}
		}

		pool, err := pants.NewManifoldStatePool(ctx, mf, &wg,
			pants.WithSize(uint(poolSize)),
			pants.WithStateInitializer(initializer),
			pants.WithStateFinalizer(finalizer),
			pants.WithOutput(uint(jobCount), time.Microsecond*100, time.Second),
		)
		Expect(err).NotTo(HaveOccurred())

		commands := []string{"echo 1", "echo 2", "echo 3", "echo 4", "echo 5"}
		for _, cmd := range commands {
			err := pool.Post(ctx, cmd)
			Expect(err).NotTo(HaveOccurred())
		}

		pool.Conclude(ctx)

		var results []string //nolint:prealloc // not necessary to prealloc for a test
		for output := range pool.Observe() {
			Expect(output.Error).NotTo(HaveOccurred())
			results = append(results, output.Payload)
		}

		Expect(len(results)).To(Equal(jobCount))
		for i, res := range results {
			Expect(res).To(ContainSubstring("executed: echo"))
			_ = i
		}

		pool.Release(ctx)
	})
})
