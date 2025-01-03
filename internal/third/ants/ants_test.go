package ants_test

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/pants"
	"github.com/snivilised/pants/internal/lab"
	"github.com/snivilised/pants/internal/third/ants"
	"github.com/snivilised/pants/locale"
)

var _ = Describe("Ants", Ordered, func() {
	var (
		repo                string
		l10nPath            string
		testTranslationFile li18ngo.TranslationFiles
	)

	BeforeAll(func() {
		repo = lab.Repo("")
		l10nPath = lab.Path(repo, "test/data/l10n")

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

	Context("NewPool", func() {
		Context("Submit", func() {
			When("non-blocking", func() {
				It("🧪 should: not fail", func(specCtx SpecContext) {
					// TestNonblockingSubmit
					// ??? defer leaktest.Check(GinkgoT())()
					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					const poolSize = 10

					pool, err := ants.NewPool(ctx,
						ants.WithSize(poolSize),
						ants.WithNonblocking(true),
					)

					Expect(err).To(Succeed(), "create TimingPool failed")

					defer pool.Release(ctx)

					for i := 0; i < poolSize-1; i++ {
						Expect(pool.Submit(ctx, longRunningFunc)).To(Succeed(),
							"nonblocking submit when pool is not full shouldn't return error",
						)
					}
					firstCh := make(chan struct{})
					secondCh := make(chan struct{})
					fn := func() {
						<-firstCh
						close(secondCh)
					}
					// p is full now.
					Expect(pool.Submit(ctx, fn)).To(Succeed(),
						"nonblocking submit when pool is not full shouldn't return error",
					)
					Expect(pool.Submit(ctx, demoFunc)).To(
						MatchError(locale.ErrPoolOverload.Error()),
						"nonblocking submit when pool is full should get an ErrPoolOverload",
					)

					// interrupt fn to get an available worker
					close(firstCh)
					<-secondCh
					Expect(pool.Submit(ctx, demoFunc)).To(Succeed(),
						"nonblocking submit when pool is not full shouldn't return error",
					)
				})
			})

			When("max blocking", func() {
				It("🧪 should: not fail", func(specCtx SpecContext) {
					// TestMaxBlockingSubmit
					// ??? defer leaktest.Check(GinkgoT())()
					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					const poolSize = 10
					pool, err := ants.NewPool(ctx,
						ants.WithSize(poolSize),
						ants.WithMaxBlockingTasks(1),
					)
					Expect(err).To(Succeed(), "create TimingPool failed")

					defer pool.Release(ctx)

					for i := 0; i < poolSize-1; i++ {
						Expect(pool.Submit(ctx, longRunningFunc)).To(Succeed(),
							"blocking submit when pool is not full shouldn't return error",
						)
					}
					ch := make(chan struct{})
					fn := func() {
						<-ch
					}
					// p is full now.
					Expect(pool.Submit(ctx, fn)).To(Succeed(),
						"nonblocking submit when pool is not full shouldn't return error",
					)

					var wg sync.WaitGroup
					wg.Add(1)
					errCh := make(chan error, 1)
					go func() {
						// should be blocked. blocking num == 1
						if err := pool.Submit(ctx, demoFunc); err != nil {
							errCh <- err
						}
						wg.Done()
					}()
					time.Sleep(1 * time.Second)
					// already reached max blocking limit
					Expect(pool.Submit(ctx, demoFunc)).To(
						MatchError(locale.ErrPoolOverload.Error()),
						"blocking submit when pool reach max blocking submit should return ErrPoolOverload",
					)

					// interrupt f to make blocking submit successful.
					close(ch)
					wg.Wait()
					select {
					case <-errCh:
						// t.Fatalf("blocking submit when pool is full should not return error")
						Fail("blocking submit when pool is full should not return error")
					default:
					}
				})
			})
		})

		Context("IfOption", func() {
			When("true", func() {
				It("🧪 should: use option", func(specCtx SpecContext) {
					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					const (
						poolSize = 10
					)

					pool, _ := ants.NewPool(ctx,
						ants.IfOption(true, ants.WithSize(poolSize)),
						ants.WithNonblocking(true),
					)

					options := pool.GetOptions()
					Expect(options.Size).To(BeEquivalentTo(poolSize))
				})
			})

			When("false", func() {
				It("🧪 should: use option", func(specCtx SpecContext) {
					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					const (
						poolSize = 10
					)

					pool, _ := ants.NewPool(ctx,
						ants.IfOption(false, ants.WithSize(poolSize)),
						ants.WithNonblocking(true),
					)

					options := pool.GetOptions()
					Expect(options.Size).To(BeEquivalentTo(runtime.NumCPU()))
				})
			})
		})
	})

	Context("NewPoolWithFunc", func() {
		Context("Invoke", func() {
			When("waiting to get worker", func() {
				It("🧪 should: not fail", func(specCtx SpecContext) {
					// TestAntsPoolWithFuncWaitToGetWorker

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					var wg sync.WaitGroup
					pool, _ := ants.NewPoolWithFunc(ctx, func(i pants.InputEnvelope) {
						demoPoolFunc(i)
						wg.Done()
					},
						ants.WithSize(AntsSize),
					)
					defer pool.Release(ctx)

					for i := 0; i < n; i++ {
						wg.Add(1)
						_ = pool.Invoke(ctx, Param)
					}
					wg.Wait()
					GinkgoWriter.Printf("pool with func, no of running workers:%d\n",
						pool.Running(),
					)
					ShowMemStats()
				})
			})

			When("waiting to get worker with pre malloc", func() {
				It("🧪 should: not fail", func(specCtx SpecContext) {
					// TestAntsPoolWithFuncWaitToGetWorkerPreMalloc

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					var wg sync.WaitGroup
					pool, _ := ants.NewPoolWithFunc(ctx, func(i pants.InputEnvelope) {
						demoPoolFunc(i)
						wg.Done()
					},
						ants.WithSize(AntsSize),
						ants.WithPreAlloc(true),
					)
					defer pool.Release(ctx)

					for i := 0; i < n; i++ {
						wg.Add(1)
						_ = pool.Invoke(ctx, Param)
					}
					wg.Wait()
					GinkgoWriter.Printf("pool with func, no of running workers:%d\n",
						pool.Running(),
					)
					ShowMemStats()
				})
			})
		})
	})
})
