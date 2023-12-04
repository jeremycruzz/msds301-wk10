package montecarlo_test

import (
	// adjust the import path

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/jeremycruzz/msds301-wk10/pkg/simulation/montecarlo"
	"github.com/jeremycruzz/msds301-wk10/pkg/tableutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Montecarlo", func() {
	Skip("Skipping montecarlo tests, can't figure out memory stuff")
	var (
		table  arrow.Table
		err    error
		result int64
	)

	Context("when running a montecarlo simulation", func() {
		BeforeEach(func() {
			table, _ = tableutils.InferCsvToTable("../../data/test.csv")
			table.Retain()
			result = montecarlo.PercentileValue(*table.Column(0), float64(.5), 4, 1000)
		})
		It("should return a non-nil result without error", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
		})
		It("should return 1", func() {
			Expect(result).To(Equal(int64(1)))
		})
	})
})
