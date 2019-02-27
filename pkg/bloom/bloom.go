package bloom

import (
	"fmt"
	"hash"

	"github.com/spaolacci/murmur3"
)

type bloom struct {
	// Public Vars
	Capacity  int
	ErrorRate float32

	// Private vars
	hashers []hash.Hash64
}

// New creates an instance of the bloom filter
func New(hashicity, capacity int, errorRate float32) (bloom, error) {
	/*
		Initialize the bloom filter here, need to do a few things
		- Initialize the hash functions?
		- Create the hash array
		- hopefully not die
	*/
	var b bloom

	hashers := make([]hash.Hash64, 0)
	for i := 0; i < hashicity; i++ {
		hashers = append(hashers, murmur3.New64WithSeed(uint32(i)))
	}

	// Preconditions:
	if capacity == 0 {
		return b, fmt.Errorf("capacity cannot be zero")
	}
	if !(errorRate > 0) || !(errorRate < 100) {
		return b, fmt.Errorf("probability must be greater than 0 and less than 100")
	}

	// Initialize the hash functions

	b = bloom{
		Capacity:  capacity,
		ErrorRate: errorRate,
	}
	return b, nil
}
