# msds301-wk10

### Setup (windows, gitbash)
- Clone repo with `git clone git@github.com:jeremycruzz/msds301-wk10.git`
- Run `go mod tidy`

### Building executable
- Run `go build -o sprintsim.exe ./cmd/sprintsimulator/main.go`
  - I'm not sure why I wasn't able to run `./cmd/sprintsimulator` this time

### Testing
- Run `go test ./...`
- Right now montecarlo tests are skipped because of memory issues with tests.

### Running Go executable
- Run `./sprintsim.exe {threads} {runMemoryProfiling}` 
    - default threads is `4`
    - to run memory profiling `runMemoryProfiling == true`
      - generate mem report `go tool pprof -pdf memprofile.pprof > profile2.pdf`

### Benchmarking
- Uncomment the memory parts in the r script
- run `./benchmark.sh 1000`

### Background / Conclusion

For this assignment I originally wanted to try to implement MonteCarloMakarovChains but I struggled to understand it well enough to implement it and decided to do Monte Carlo simulations. Since the Monte Carlo simulations were pretty straight forward to implement I wanted to try doing it with Apache Arrow which was not straight forward to use. I also ended up using mostly the base R package to implement the Simulations. After implementing both we can see similar results:

<details> 
<Summary> Results </Summary>

Results were slightly different because of the nature of random sampling

###### Go
Jeremy's 50th percentile: 5
Angela's 50th percentile: 6
Mike's 50th percentile: 5
Team's 50th percentile: 17
Average number of sprints to complete 100 pts: 6

###### R
Jeremy's 50th percentile: 6 
Angelica's 50th percentile: 7
Mike's 50th percentile: 5
Team's 50th percentile: 18
Average number of sprints to complete 100 pts: 6.187

</details>

In order to improve on my implmented version of the monte carlo simulations, I added support for multi-threading and used Apache Arrow as my data scructure in Go. Using Apache Arrow wasn't so straight forward. In my function `InferCsvToTable` I ran into issues where my table wasn't being populated with the records. After a bit of digging I found out that `NewTableFromRecords` released the columns of a record and I had to manually `Retain()` each record in order to put it in my table. While its not intuitive, it makes sense once you understand how Table columns are created in arrow. Table columns are actually `chunks` which is a list of columns that could be disjointed in memory that can be used as one column. So instead of recreating the column in a sequence of memory, the table structure just claims each record table as its own. Here are the results over 1000 runs:

|                  | R           | Go(1)      | Go(2)      | Go(4)      | Go(8)      |
|------------------|-------------|------------|------------|------------|------------|
| Average Time (ns)| 526584043   | 29387303   | 29745959   | 29704564   | 30257033   |
| Runtime Allocation(kB) | 18700.00   | 1124.13  | 1124.13  | 1171.00  | 609.50   |

I was surprised that multi-threading ended up slowing down the run time in go but I think with a larger data set we might see the reverse happen. Apache Arrow's memory efficency is pretty evident by the fact that memory usage didn't increase as threads increased and even dropped by half when using 8 treads (even after multiple runs.) Comparing Go to R we can see that the slowest run (Go(8)) was 31x faster than R. Additionally we can see that memory usage was 15-30x more in R than in Go. 

With these results its clear that Go is much faster and memory efficent than R. As our data grows, this will become more valuable. However, implementing some of these statistical functions in Go is a bit harder than implmenting them in R and R has more community support. As a data scientist I think we can try to create an open source community to help us create these statistical methods in Go to offset the disadvantage of using Go. 

### Possible Extension
- We could implement Apache Arrow in R and see if the time/memory usages go down.