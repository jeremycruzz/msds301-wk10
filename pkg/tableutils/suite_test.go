package tableutils_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTableutils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tableutils Suite")
}
