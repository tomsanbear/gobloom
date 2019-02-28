package bloom

import (
	"fmt"
	"strconv"
	"testing"
	"time"

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
	sut, err := New(100, 0.001)
	assert.NoError(t, err)

	testCycles := 100

	// Start the test
	startTime := time.Now().UnixNano()

	// Load every other entry
	for i := 0; i < testCycles; i = i + 2 {
		iAsString := strconv.Itoa(i) + "meow"
		currentTime := time.Now().UnixNano()
		err := sut.Add(iAsString)
		currentTime2 := time.Now().UnixNano()
		assert.NoError(t, err)
		fmt.Printf("iter: %v, add item: \"%v\" took: %v ns\n", i, iAsString, currentTime2-currentTime)
	}

	for i := 0; i < testCycles; i++ {
		iAsString := strconv.Itoa(i) + "meow"
		currentTime := time.Now().UnixNano()
		ok, err := sut.Check(iAsString)
		currentTime2 := time.Now().UnixNano()
		assert.NoError(t, err)
		fmt.Printf("iter: %v, get item: \"%v\", result: %v took: %v ns\n", i, iAsString, ok, currentTime2-currentTime)
	}

	// End the test
	endTime := time.Now().UnixNano()

	fmt.Printf("Test took: %v ns\n", endTime-startTime)
}
