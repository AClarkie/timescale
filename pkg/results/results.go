package results

import (
	"fmt"
	"sort"
	"time"
)

type Result struct {
	QueryTimes []time.Duration
	Error      error
}

func ProcessAndDisplayResults(results chan *Result) {

	sortedResults := aggregateQueryTimesAndSort(results)
	totalQueries := len(sortedResults)
	totalTime := calculateTotalTime(sortedResults)

	fmt.Println("RESULTS")
	fmt.Println("=======")
	fmt.Printf("Number of queries run:                    %d \n", totalQueries)
	fmt.Printf("Total processing time across all queries: %v \n", totalTime)
	fmt.Printf("The average query time:                   %v \n", totalTime/time.Duration(totalQueries))
	fmt.Printf("The minimum query time:                   %v \n", sortedResults[0])
	fmt.Printf("The maximum query time:                   %v \n", sortedResults[totalQueries-1])
	fmt.Printf("The median query time:                    %v \n", calculateMedianTime(sortedResults))
}

func aggregateQueryTimesAndSort(results chan *Result) []time.Duration {
	var queryTimes []time.Duration
	for s := range results {
		queryTimes = append(queryTimes, s.QueryTimes...)
	}
	sort.Slice(queryTimes, func(i, j int) bool { return queryTimes[i] < queryTimes[j] })
	return queryTimes
}

func calculateMedianTime(queryTimes []time.Duration) time.Duration {
	length := len(queryTimes)
	if length%2 == 0 {
		return (queryTimes[length/2] + queryTimes[length/2-1]) / 2
	} else {
		return queryTimes[length/2]
	}
}

func calculateTotalTime(queryTimes []time.Duration) time.Duration {
	var totalTime time.Duration
	for _, time := range queryTimes {
		totalTime += time
	}
	return totalTime
}
