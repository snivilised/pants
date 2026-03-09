package locale_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLocale(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "locale Suite")
}
