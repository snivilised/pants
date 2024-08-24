package pants_test

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/pants"
	"github.com/snivilised/pants/internal/ants"
	"github.com/snivilised/pants/internal/helpers"
)

func produce(ctx context.Context,
	pool *pants.ManifoldFuncPool[int, int],
	wg pants.WaitGroup,
) {
	defer wg.Done()

	for i, n := 0, 100; i < n; i++ {
		_ = pool.Post(ctx, Param)
	}

	pool.Conclude(ctx)
}

func inject(ctx context.Context,
	pool *pants.ManifoldFuncPool[int, int],
	wg pants.WaitGroup,
) {
	defer wg.Done()

	ch := pool.Source(ctx, wg)
	for i, n := 0, 100; i < n; i++ {
		ch <- Param
	}

	close(ch)
}

func consume(_ context.Context,
	pool *pants.ManifoldFuncPool[int, int],
	wg pants.WaitGroup,
) {
	defer wg.Done()

	for output := range pool.Observe() {
		Expect(output.Error).To(Succeed())
		Expect(output.ID).NotTo(BeEmpty())
		Expect(output.SequenceNo).NotTo(Equal(0))
	}
}

var _ = Describe("WorkerPoolFuncManifold", Ordered, func() {
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
			Context("with consumer", func() {
				It("🧪 should: not fail", func(specCtx SpecContext) {
					// TestNonblockingSubmit
					var wg sync.WaitGroup

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					pool, err := pants.NewManifoldFuncPool(
						ctx, demoPoolManifoldFunc, &wg,
						pants.WithSize(AntsSize),
						pants.WithOutput(10, CheckCloseInterval, TimeoutOnSend),
					)

					defer pool.Release(ctx)

					wg.Add(1)
					go produce(ctx, pool, &wg)

					wg.Add(1)
					go consume(ctx, pool, &wg)

					wg.Wait()
					GinkgoWriter.Printf("pool with func, no of running workers:%d\n",
						pool.Running(),
					)
					ShowMemStats()

					Expect(err).To(Succeed())
				})
			})

			Context("without consumer", func() {
				It("🧪 should: not fail", func(specCtx SpecContext) {
					// TestNonblockingSubmit
					var wg sync.WaitGroup

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					pool, err := pants.NewManifoldFuncPool(
						ctx, demoPoolManifoldFunc, &wg,
						pants.WithSize(AntsSize),
					)

					defer pool.Release(ctx)

					wg.Add(1)
					go produce(ctx, pool, &wg)

					wg.Wait()
					GinkgoWriter.Printf("pool with func, no of running workers:%d\n",
						pool.Running(),
					)
					ShowMemStats()

					Expect(err).To(Succeed())
				})
			})

			Context("with input stream", func() {
				It("🧪 should: not fail", func(specCtx SpecContext) {
					// TestNonblockingSubmit
					var wg sync.WaitGroup

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					pool, err := pants.NewManifoldFuncPool(
						ctx, demoPoolManifoldFunc, &wg,
						pants.WithSize(AntsSize),
						pants.WithInput(InputBufferSize),
						pants.WithOutput(10, CheckCloseInterval, TimeoutOnSend),
					)

					defer pool.Release(ctx)

					wg.Add(1)
					go inject(ctx, pool, &wg)

					wg.Add(1)
					go consume(ctx, pool, &wg)

					wg.Wait()
					GinkgoWriter.Printf("pool with func, no of running workers:%d\n",
						pool.Running(),
					)
					ShowMemStats()

					Expect(err).To(Succeed())
				})
			})

			Context("cancelled", func() {
				Context("without consumer", func() {
					It("🧪 should: not fail", func(specCtx SpecContext) {
						// TestNonblockingSubmit
						var wg sync.WaitGroup

						ctx, cancel := context.WithCancel(specCtx)
						defer cancel()

						pool, err := pants.NewManifoldFuncPool(
							ctx, demoPoolManifoldFunc, &wg,
							pants.WithSize(AntsSize),
						)

						defer pool.Release(ctx)

						wg.Add(1)
						go func(ctx context.Context,
							pool *pants.ManifoldFuncPool[int, int],
							wg pants.WaitGroup,
						) {
							defer wg.Done()

							for i, n := 0, 100; i < n; i++ {
								_ = pool.Post(ctx, Param)

								if i > 10 {
									cancel()
									break
								}
							}
							pool.Conclude(ctx)
						}(ctx, pool, &wg)

						wg.Wait()
						GinkgoWriter.Printf("pool with func, no of running workers:%d\n",
							pool.Running(),
						)
						ShowMemStats()

						Expect(err).To(Succeed())
					})
				})
			})

			Context("timeout on send, with cancellation monitor", func() {
				When("output requested, but accidentally not consumed by client", func() {
					It("🧪 should: cancel context and terminate", func(specCtx SpecContext) {
						// TestNonblockingSubmit
						var wg sync.WaitGroup

						ctx, cancel := context.WithCancel(specCtx)
						defer cancel()

						pool, err := pants.NewManifoldFuncPool(
							ctx, demoPoolManifoldFunc, &wg,
							pants.WithSize(AntsSize),
							pants.WithInput(InputBufferSize),
							pants.WithOutput(10, CheckCloseInterval, TimeoutOnSend),
						)

						defer pool.Release(ctx)

						wg.Add(1)
						go inject(ctx, pool, &wg)

						pants.StartCancellationMonitor(ctx,
							cancel,
							&wg,
							pool.CancelCh(),
							func() {},
						)
						wg.Wait()
						GinkgoWriter.Printf("pool with func, no of running workers:%d\n",
							pool.Running(),
						)
						ShowMemStats()

						Expect(err).To(Succeed())
					})
				})
			})
		})

		Context("IfOption", func() {
			When("true", func() {
				It("🧪 should: use option", func(specCtx SpecContext) {
					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					var wg sync.WaitGroup

					const (
						poolSize = 10
					)

					pool, _ := pants.NewManifoldFuncPool(
						ctx, demoPoolManifoldFunc, &wg,
						ants.IfOption(true, ants.WithSize(poolSize)),
						pants.WithInput(InputBufferSize),
						pants.WithOutput(10, CheckCloseInterval, TimeoutOnSend),
					)

					options := pool.GetOptions()
					Expect(options.Size).To(BeEquivalentTo(poolSize))
				})
			})

			When("false", func() {
				It("🧪 should: use option", func(specCtx SpecContext) {
					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					var wg sync.WaitGroup

					const (
						poolSize = 10
					)

					pool, _ := pants.NewManifoldFuncPool(
						ctx, demoPoolManifoldFunc, &wg,
						ants.IfOption(false, ants.WithSize(poolSize)),
						pants.WithInput(InputBufferSize),
						pants.WithOutput(10, CheckCloseInterval, TimeoutOnSend),
					)

					options := pool.GetOptions()
					Expect(options.Size).To(BeEquivalentTo(runtime.NumCPU()))
				})
			})
		})
	})
})
