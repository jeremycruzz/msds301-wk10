package tableutils

import (
	"fmt"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
)

// SumCols sums the values of two or more columns and returns a new column with the sum of each row
// For now there is only support for 1 array in a chunk
func SumCols(cols ...*arrow.Column) (arrow.Column, error) {
	if len(cols) == 0 {
		return arrow.Column{}, fmt.Errorf("no columns specified: two or more columns required")
	}

	if len(cols) == 1 {
		return arrow.Column{}, fmt.Errorf("only one column specified: two or more columns required")
	}

	//create new int64 array
	sumCol := array.NewInt64Builder(memory.DefaultAllocator)

	//sum the total of each row
	for rows := 0; rows < cols[0].Len(); rows++ {
		sum := int64(0)
		for _, column := range cols {
			sum += column.Data().Chunks()[0].(*array.Int64).Value(rows) //expecting there to only be 1 chunk for all data
		}
		sumCol.Append(sum)
	}

	field := arrow.Field{Name: "Sum", Type: arrow.PrimitiveTypes.Int64, Nullable: true}
	newCol := arrow.NewColumnFromArr(field, sumCol.NewInt64Array())

	return newCol, nil
}
