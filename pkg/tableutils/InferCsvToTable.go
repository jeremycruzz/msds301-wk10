package tableutils

import (
	"fmt"
	"io"
	"os"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/csv"
	"github.com/apache/arrow/go/v14/arrow/memory"
)

func InferCsvToTable(csvPath string) (arrow.Table, error) {

	//open file
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV: %w", err)
	}
	defer f.Close()

	// read csv file
	reader := csv.NewInferringReader(f,
		csv.WithHeader(true),
		csv.WithAllocator(memory.DefaultAllocator),
		csv.WithChunk(-1), // -1 means read all rows at once
		csv.WithNullReader(false))
	defer reader.Release()

	// add each row to the records array
	var records []arrow.Record
	for reader.Next() {
		rec := reader.Record()
		rec.Retain() // retain the record so it doesn't get GC'd
		records = append(records, rec)
	}

	if reader.Err() != nil && reader.Err() != io.EOF {
		return nil, fmt.Errorf("error reading CSV: %w", reader.Err())
	}

	// create table
	table := array.NewTableFromRecords(records[0].Schema(), records) // this releases the records and is the reason rec.Retain() is needed above
	return table, nil
}
