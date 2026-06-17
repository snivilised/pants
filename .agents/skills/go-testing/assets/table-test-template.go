package example_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Example", func() {
	DescribeTable("behaviour",
		func(give string, want string) {
			// TODO: call function under test
			got := give // TODO: replace with actual call
			Expect(got).To(Equal(want))
		},
		Entry("basic case", "", ""),
		Entry("edge case", "", ""),
	)
})
