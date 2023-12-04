package montecarlo

import (
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
)

func PercentileValue(column arrow.Column, percentile float64, threads, samples int) int64 {
	var wg sync.WaitGroup
	split := samples / threads
	remainder := samples % threads

	var lock sync.Mutex
	sharedFreqMap := make(map[int64]int)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(thread int) {
			column.Retain()
			defer column.Release()
			defer wg.Done()
			localFreqMap := make(map[int64]int)
			localRand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))) //hopfully this is unique enough

			// get samples for this thread
			localSamples := split
			if thread < remainder {
				localSamples++ //take care of remainder
			}

			// get random values
			for j := 0; j < localSamples; j++ {
				randomIndex := localRand.Intn(int(column.Len()))
				value := (column.Data().Chunks()[0].(*array.Int64).Value(randomIndex))
				localFreqMap[value]++
			}

			//add values to shared map
			lock.Lock()
			for value, count := range localFreqMap {
				sharedFreqMap[value] += count
			}
			lock.Unlock()
		}(i)

	}
	wg.Wait()
	return calculatePercentile(sharedFreqMap, percentile)
}

func SimulateIterations(column arrow.Column, percentile float64, threads, samples, goal int) int64 {
	var wg sync.WaitGroup
	split := samples / threads
	remainder := samples % threads

	var lock sync.Mutex
	sharedFreqMap := make(map[int64]int)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(thread int) {
			column.Retain()
			defer column.Release()
			defer wg.Done()
			localFreqMap := make(map[int64]int)
			localRand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))) //hopfully this is unique enough

			// get samples for this thread
			localSamples := split
			if thread < remainder {
				localSamples++ //take care of remainder
			}

			// get random values
			for j := 0; j < localSamples; j++ {
				score := int64(0)
				iterations := int64(0)
				for score < int64(goal) {
					iterations++
					randomIndex := localRand.Intn(int(column.Len()))
					score += (column.Data().Chunks()[0].(*array.Int64).Value(randomIndex))
				}

				localFreqMap[iterations]++
			}

			//add values to shared map
			lock.Lock()
			for value, count := range localFreqMap {
				sharedFreqMap[value] += count
			}
			lock.Unlock()
		}(i)

	}
	wg.Wait()
	return calculatePercentile(sharedFreqMap, percentile)
}

func calculatePercentile(freqMap map[int64]int, percentile float64) int64 {

	totalSamples := 0
	for _, count := range freqMap {
		totalSamples += count
	}

	index := int(math.Round(float64(totalSamples) * percentile))

	//sort keys in ascending order
	keys := make([]int64, 0, len(freqMap))
	for k := range freqMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	// find the value at the index
	count := 0
	for _, k := range keys {
		count += freqMap[k]
		if count >= index {
			return k
		}
	}

	return keys[len(keys)-1] // 100% case
}
