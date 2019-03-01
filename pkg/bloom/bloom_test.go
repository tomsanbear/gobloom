package bloom

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tables := []struct {
		capacity      int
		errorRate     float64
		expectedError bool
	}{
		{1, 0.1, false},
		{0, 0.1, true},
		{1, 100.01, true},
	}

	for _, table := range tables {
		_, err := New(table.capacity, table.errorRate)

		if table.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestBitArraySize(t *testing.T) {
	tables := []struct {
		n float64
		p float64
		m float64
	}{
		{20, 0.05, 124},
	}

	for _, table := range tables {
		actualM := optimalBitArraySize(table.n, table.p)
		assert.EqualValues(t, table.m, actualM)
	}
}

func TestOptimalHashFunk(t *testing.T) {
	tables := []struct {
		m        int
		n        int
		numFuncs int
	}{
		{124, 20, 4},
	}
	for _, table := range tables {
		actualM := optimalHashFunk(table.m, table.n)
		assert.EqualValues(t, table.numFuncs, actualM)
	}
}

func TestBloom(t *testing.T) {
	tables := []struct {
		testCycles   int
		predictedCap int
		errorRate    float64
	}{
		{10, 10, 0.01},
		{100, 100, 0.01},
		{1000, 1000, 0.01},
		{10000, 10000, 0.01},
		{1000000, 1000000, 0.01},
	}

	for _, table := range tables {
		// Setup for tracking performance across test runs
		var putTimes []int64
		var getTimes []int64

		sut, err := New(table.predictedCap, table.errorRate)
		assert.NoError(t, err)

		// Start the test
		startTime := time.Now().UnixNano()

		// Load every other entry
		for i := 0; i < table.testCycles; i = i + 2 {
			iAsString := strconv.Itoa(i) + ".meow.com"
			currentTime := time.Now().UnixNano()
			err := sut.Add(iAsString)
			currentTime2 := time.Now().UnixNano()
			assert.NoError(t, err)
			putTimes = append(putTimes, currentTime2-currentTime)
		}

		falsePositives := 0
		for i := 0; i < table.testCycles; i++ {
			iAsString := strconv.Itoa(i) + ".meow.com"
			currentTime := time.Now().UnixNano()
			ok, err := sut.Check(iAsString)
			currentTime2 := time.Now().UnixNano()
			if (i % 2) != 0 {
				if ok == true {
					falsePositives++
				}
			}
			assert.NoError(t, err)
			getTimes = append(getTimes, currentTime2-currentTime)
		}

		// Find the interesting stats for the test run
		putStats := stats.LoadRawData(putTimes)
		getStats := stats.LoadRawData(getTimes)

		putMedian, err := putStats.Median()
		assert.NoError(t, err)

		getMedian, err := getStats.Median()
		assert.NoError(t, err)

		putVar, err := putStats.Variance()
		assert.NoError(t, err)

		getVar, err := getStats.Variance()
		assert.NoError(t, err)

		// End the test
		endTime := time.Now().UnixNano()

		// Reporting data
		fmt.Printf("Test took: %v ns\n", endTime-startTime)
		fmt.Printf("Setup: testCycles = %v, capacity = %v, desiredErrorRate = %v\n", table.testCycles, table.predictedCap, table.errorRate)
		fmt.Printf("False positive rate was: %v\nPUT median: %v ns, variance: %v\n GET median %v ns, variance: %v\n", float64(falsePositives)/float64(table.testCycles), putMedian, putVar, getMedian, getVar)
	}

}
