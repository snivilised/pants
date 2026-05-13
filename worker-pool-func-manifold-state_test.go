package pants_test

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/pants"
)

type mockState struct {
	id    pants.RoutineID
	count int32
}

var _ = Describe("ManifoldStatePool", func() {
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

	It("should initialize state exactly once per worker and pass it to jobs", func() {
		var initCount int32
		var finalCount int32

		poolSize := 2
		jobCount := 10

		initializer := func(id pants.RoutineID) interface{} {
			atomic.AddInt32(&initCount, 1)
			return &mockState{id: id}
		}

		finalizer := func(state interface{}) {
			atomic.AddInt32(&finalCount, 1)
		}

		mf := func(input int, state *mockState) (int, error) {
			atomic.AddInt32(&state.count, 1)
			return input * 2, nil
		}

		pool, err := pants.NewManifoldStatePool(ctx, mf, &wg,
			pants.WithSize(uint(poolSize)),
			pants.WithStateInitializer(initializer),
			pants.WithStateFinalizer(finalizer),
			pants.WithOutput(uint(jobCount), time.Microsecond*100, time.Second),
		)
		Expect(err).NotTo(HaveOccurred())

		for i := 0; i < jobCount; i++ {
			err := pool.Post(ctx, i)
			Expect(err).NotTo(HaveOccurred())
		}

		pool.Conclude(ctx)

		var totalJobsHandled int32
		for output := range pool.Observe() {
			Expect(output.Error).NotTo(HaveOccurred())
			totalJobsHandled++
		}

		Expect(totalJobsHandled).To(Equal(int32(jobCount)))
		Expect(atomic.LoadInt32(&initCount)).To(BeNumerically("<=", poolSize))

		// Release the pool to trigger finalizers
		pool.Release(ctx)
		
		// The number of finalizers should match the number of initializers
		Eventually(func() int32 {
			return atomic.LoadInt32(&finalCount)
		}).Should(Equal(atomic.LoadInt32(&initCount)))
	})
})
