package pants_test

import (
	"fmt"
	"os"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/pants"
	"github.com/snivilised/pants/internal/lab"
)

var _ = Describe("pants.WaitGroup", Ordered, func() {
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
