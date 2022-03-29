package results

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestCalculateMedianTime(t *testing.T) {
	l := []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second, 5 * time.Second}
	result := calculateMedianTime(l)
	if result != 3*time.Second {
		t.Errorf("calculateMedianTime was incorrect, got: %v, want: %v.", result, 3*time.Second)
	}
}

func TestCalculateTotalTime(t *testing.T) {
	l := []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second, 5 * time.Second}
	result := calculateTotalTime(l)
	if result != 15*time.Second {
		t.Errorf("calculateTotalTime was incorrect, got: %v, want: %v.", result, 3*time.Second)
	}
}

func TestAggregateQueryTimesAndSort(t *testing.T) {
	resultsChan := make(chan *Result)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		resultsChan <- &Result{QueryTimes: []time.Duration{6 * time.Second, 2 * time.Second, 8 * time.Second, 4 * time.Second, 5 * time.Second}, Error: nil}
		resultsChan <- &Result{QueryTimes: []time.Duration{1 * time.Second, 7 * time.Second, 3 * time.Second, 9 * time.Second, 10 * time.Second}, Error: nil}
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	result := aggregateQueryTimesAndSort(resultsChan)
	expectedResult := []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second, 5 * time.Second, 6 * time.Second, 7 * time.Second, 8 * time.Second, 9 * time.Second, 10 * time.Second}
	if !compareTimeSlicesExactly(result, expectedResult) {
		t.Errorf("aggregateQueryTimesAndSort was incorrect, got: %v, want: %v.", result, expectedResult)
	}
}

func TestAggregateQueryTimesAndSortWithError(t *testing.T) {
	resultsChan := make(chan *Result)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		resultsChan <- &Result{QueryTimes: []time.Duration{6 * time.Second, 2 * time.Second, 8 * time.Second, 4 * time.Second, 5 * time.Second}, Error: nil}
		resultsChan <- &Result{QueryTimes: nil, Error: errors.New("Query failed for hostname")}
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	result := aggregateQueryTimesAndSort(resultsChan)
	expectedResult := []time.Duration{2 * time.Second, 4 * time.Second, 5 * time.Second, 6 * time.Second, 8 * time.Second}
	if !compareTimeSlicesExactly(result, expectedResult) {
		t.Errorf("aggregateQueryTimesAndSort was incorrect, got: %v, want: %v.", result, expectedResult)
	}
}

func compareTimeSlicesExactly(s1, s2 []time.Duration) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
