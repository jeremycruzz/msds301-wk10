package montecarlo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMontecarlo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Montecarlo Suite")
}
