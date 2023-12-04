#!/bin/bash

# Check if an argument is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <number of runs>"
    exit 1
fi

# Number of times the programs should run
num_runs=$1

total_time_go=0
total_time_r=0

for i in $(seq 1 $num_runs)
do
   # Go
   start_time_go=$(date +%s%N)
   ./sprintsim.exe 4
   end_time_go=$(date +%s%N)
   total_time_go=$((total_time_go + end_time_go - start_time_go))

   # R
   start_time_r=$(date +%s%N)
   Rscript ./r/sprintsimi.r
   end_time_r=$(date +%s%N)
   total_time_r=$((total_time_r + end_time_r - start_time_r))
done

# Calculate the averages
avg_time_go=$((total_time_go / num_runs))
avg_time_r=$((total_time_r / num_runs))

echo "Average time for Go: $avg_time_go nanoseconds"
echo "Average time for R: $avg_time_r nanoseconds"