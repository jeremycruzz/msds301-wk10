package tableutils_test

import (
	// adjust the import path
	"fmt"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/jeremycruzz/msds301-wk10/pkg/tableutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InferCsvToTable", func() {
	var (
		table arrow.Table
		err   error
	)
	Context("when reading a CSV file with a header", func() {
		BeforeEach(func() {
			table, err = tableutils.InferCsvToTable("../../data/test.csv")
		})
		It("should read a CSV file and return a non-nil table without error", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(table).NotTo(BeNil())
		})
		It("should have 3 columns", func() {
			Expect(table.NumCols()).To(Equal(int64(3)))
		})
		It("have all the correct data", func() {
			Expect(fmt.Sprint(table)).To(ContainSubstring("[[1 1 1 1]]"))
			Expect(fmt.Sprint(table)).To(ContainSubstring("[[2 2 2 2]]"))
			Expect(fmt.Sprint(table)).To(ContainSubstring("[[5 5 5 5]]"))
		})
	})

	Context("when reading a non existent file", func() {
		BeforeEach(func() {
			table, err = tableutils.InferCsvToTable("none.csv")
		})

		It("should return an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})
})
