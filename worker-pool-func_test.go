package pants_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/pants"
	"github.com/snivilised/pants/internal/helpers"
	"github.com/snivilised/pants/locale"
)

var _ = Describe("WorkerPoolFunc", Ordered, func() {
	var (
		repo                string
		l10nPath            string
		testTranslationFile li18ngo.TranslationFiles
	)

	BeforeAll(func() {
		repo = helpers.Repo("")
		l10nPath = helpers.Path(repo, "test/data/l10n")

		_, err := os.Stat(l10nPath)
		Expect(err).To(Succeed(),
			fmt.Sprintf("l10n '%v' path does not exist", l10nPath),
		)

		testTranslationFile = li18ngo.TranslationFiles{
			li18ngo.Li18ngoSourceID: li18ngo.TranslationSource{Name: "test"},
		}
	})

	BeforeEach(func() {
		if err := li18ngo.Use(func(o *li18ngo.UseOptions) {
			o.Tag = li18ngo.DefaultLanguage
			o.From.Sources = testTranslationFile
		}); err != nil {
			Fail(err.Error())
		}
	})

	Context("ants", func() {
		When("NonBlocking", func() {
			It("should: not fail", func(specCtx SpecContext) {
				// TestNonblockingSubmit
				var wg sync.WaitGroup

				ctx, cancel := context.WithCancel(specCtx)
				defer cancel()

				pool, err := pants.NewFuncPool[int, int](ctx, demoPoolFunc, &wg,
					pants.WithSize(AntsSize),
				)

				defer pool.Release(ctx)

				for i := 0; i < n; i++ {
					_ = pool.Post(ctx, Param)
				}
				wg.Wait()
				GinkgoWriter.Printf("pool with func, no of running workers:%d\n",
					pool.Running(),
				)
				ShowMemStats()

				Expect(err).To(Succeed())
			})

			Context("cancelled", func() {
				It("should: not fail", func(specCtx SpecContext) {
					// TestNonblockingSubmit
					var wg sync.WaitGroup

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					pool, err := pants.NewFuncPool[int, int](ctx, demoPoolFunc, &wg,
						pants.WithSize(AntsSize),
					)

					defer pool.Release(ctx)

					for i := 0; i < n; i++ {
						_ = pool.Post(ctx, Param)

						if i > 10 {
							cancel()
							break
						}
					}
					wg.Wait()
					GinkgoWriter.Printf("pool with func, no of running workers:%d\n",
						pool.Running(),
					)
					ShowMemStats()

					Expect(err).To(Succeed())
				})
			})
		})

		When("MaxNonBlocking", func() {
			It("should: not fail", func(specCtx SpecContext) {
				// TestMaxBlockingSubmitWithFunc
				var wg sync.WaitGroup

				ctx, cancel := context.WithCancel(specCtx)
				defer cancel()

				pool, err := pants.NewFuncPool[int, int](ctx, longRunningPoolFunc, &wg,
					pants.WithSize(PoolSize),
					pants.WithMaxBlockingTasks(1),
				)

				Expect(err).To(Succeed(), "create TimingPool failed")
				defer pool.Release(ctx)

				By("👾 POOL-CREATED\n")
				for i := 0; i < PoolSize-1; i++ {
					Expect(pool.Post(ctx, Param)).To(Succeed(),
						"submit when pool is not full shouldn't return error",
					)
				}

				ch := make(chan struct{})
				// pool is full now.
				Expect(pool.Post(ctx, ch)).To(Succeed(),
					"submit when pool is not full shouldn't return error",
				)

				By("👾 WAIT-GROUP ADD(worker-pool-func)\n")
				wg.Add(1)
				errCh := make(chan error, 1)

				go func() {
					// should be blocked. blocking num == 1
					if err := pool.Post(ctx, Param); err != nil {
						errCh <- err
					}
					By("👾 Producer complete\n")
					wg.Done()
				}()
				time.Sleep(1 * time.Second)
				// already reached max blocking limit
				Expect(pool.Post(ctx, Param)).To(
					MatchError(locale.ErrPoolOverload.Error()),
					"blocking submit when pool reach max blocking submit should return ErrPoolOverload",
				)

				By("👾 CLOSING\n")
				// interrupt one func to make blocking submit successful.

				close(ch)
				wg.Wait()
				select {
				case <-errCh:
					Fail("blocking submit when pool is full should not return error")
				default:
				}
			})
		})
	})
})
