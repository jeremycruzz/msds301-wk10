package main

import (
	"fmt"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/jeremycruzz/msds301-wk10/pkg/tableutils"
)

func main() {
	var table arrow.Table
	table, err := tableutils.InferCsvToTable("./data/sprints.csv")
	table.Retain()
	defer table.Release()

	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	newCol, err := tableutils.SumCols(table.Column(1), table.Column(2), table.Column(3))
	defer newCol.Release()
	if err != nil {
		fmt.Println("Error summing columns:", err)
		return
	}

	table, err = table.AddColumn(int(table.NumCols()), newCol.Field(), newCol)
	defer table.Release()
	if err != nil {
		fmt.Println("Error adding column:", err)
		return
	}

	fmt.Println(table)
}
