package pants_test

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok

	"github.com/snivilised/pants"
)

var _ = Describe("pants.WaitGroup", func() {
	Context("given: a sync.WaitGroup", func() {
		It("should: track invocations", func() {
			var wg sync.WaitGroup
			tracker := pants.TrackWaitGroup(&wg,
				func(count int32) {
					GinkgoWriter.Printf("---> 🔴 Add (%v)\n", count)
				},
				func(count int32) {
					GinkgoWriter.Printf("---> 🟢 Done (%v)\n", count)
				},
			)

			for range 10 {
				tracker.Add(1)
				go func(tracker pants.WaitGroup) {
					defer tracker.Done()

					const delay = time.Millisecond * 100
					time.Sleep(delay)
				}(tracker)
			}

			tracker.Wait()

			if wg, ok := tracker.(*pants.TrackableWaitGroup); ok {
				Expect(wg.Count()).To(BeEquivalentTo(0))
			} else {
				Fail("tracker should be *pants.TrackableWaitGroup")
			}
		})
	})
})
