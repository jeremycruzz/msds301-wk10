package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/jeremycruzz/msds301-wk10/pkg/simulation/montecarlo"
	"github.com/jeremycruzz/msds301-wk10/pkg/tableutils"
)

func main() {
	var f *os.File
	mem := false
	threads := 4
	var err error

	args := os.Args[1:]

	if len(args) >= 1 {
		threads, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Bad Argument", err)
			return
		}
	}

	if len(args) >= 2 {
		mem = args[1] == "true"
	}

	// Memory profiling
	if mem {
		fmt.Println("Memory profiling enabled")
		f, err = os.Create("memprofile.pprof")
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
	}

	var table arrow.Table
	table, err = tableutils.InferCsvToTable("./data/sprints.csv")
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

	fmt.Println()
	fmt.Println("Given the last 28 sprints...")
	fmt.Println()
	fmt.Println("What is the mean amount of points completed by the team?")
	jer50 := montecarlo.PercentileValue(*table.Column(1), 0.5, threads, 1000)
	ang50 := montecarlo.PercentileValue(*table.Column(2), 0.5, threads, 1000)
	mik50 := montecarlo.PercentileValue(*table.Column(3), 0.5, threads, 1000)
	sum50 := montecarlo.PercentileValue(*table.Column(4), 0.5, threads, 1000)
	sum100pts := montecarlo.SimulateIterations(*table.Column(4), 0.5, threads, 1000, 100)

	fmt.Println("Jeremy's 50th percentile:", jer50)
	fmt.Println("Angela's 50th percentile:", ang50)
	fmt.Println("Mike's 50th percentile:", mik50)
	fmt.Println("Team's 50th percentile:", sum50)
	fmt.Println()
	fmt.Println("If the team needs to complete a 100 point epic, what is the average number of sprints it will take?")
	fmt.Println()
	fmt.Println("Average number of sprints to complete 100 pts:", sum100pts)

	if mem {
		fmt.Println("Writing memory profile")
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
