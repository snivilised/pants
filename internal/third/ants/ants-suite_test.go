package ants_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAnts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ants Suite")
}
