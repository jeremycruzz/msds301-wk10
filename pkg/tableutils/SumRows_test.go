package tableutils_test

import (
	// adjust the import path
	"fmt"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/jeremycruzz/msds301-wk10/pkg/tableutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SumCols", func() {
	var (
		table arrow.Table
		col   arrow.Column
		err   error
	)
	Context("when summing cols", func() {
		BeforeEach(func() {
			table, _ = tableutils.InferCsvToTable("../../data/test.csv")
			col, err = tableutils.SumCols(table.Column(0), table.Column(1), table.Column(2))
			table, err = table.AddColumn(int(table.NumCols()), col.Field(), col)
		})
		It("should read a CSV file and return a non-nil table without error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("should have 4 columns", func() {
			Expect(table.NumCols()).To(Equal(int64(4)))
		})
		It("have all the correct data", func() {
			Expect(fmt.Sprint(table)).To(ContainSubstring("[[1 1 1 1]]"))
			Expect(fmt.Sprint(table)).To(ContainSubstring("[[2 2 2 2]]"))
			Expect(fmt.Sprint(table)).To(ContainSubstring("[[5 5 5 5]]"))
			Expect(fmt.Sprint(table)).To(ContainSubstring("[[8 8 8 8]]"))
		})
	})
})
