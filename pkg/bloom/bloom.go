package bloom

import (
	"fmt"
	"math"

	"github.com/damnever/bitarray"

	"github.com/spaolacci/murmur3"
)

// Bloom interface encapsulates our useful features
type Bloom interface {
	Add(item string) error
	Check(item string) (bool, error)
}

type bloom struct {
	// Public Vars
	Capacity  int
	ErrorRate float64

	// Private vars
	numHashes    int
	hashFuncs    []murmur3.Hash128
	bitArraySize int
	bitArray     *bitarray.BitArray
}

// New creates an instance of the bloom filter
func New(capacity int, errorRate float64) (bloom, error) {
	/*
		Initialize the bloom filter here, need to do a few things
		- Initialize the hash functions
		- Create the hash array
		- hopefully not die
	*/
	var b bloom

	// Preconditions:
	if capacity == 0 {
		return b, fmt.Errorf("capacity cannot be zero")
	}
	if !(errorRate > 0) || !(errorRate < 100) {
		return b, fmt.Errorf("probability must be greater than 0 and less than 100")
	}

	// Initialize the hash functions
	baSize := optimalBitArraySize(float64(capacity), errorRate)
	nHashes := optimalHashFunk(baSize, capacity)

	hashers := make([]murmur3.Hash128, 0)
	for i := 0; i < nHashes; i++ {
		hashers = append(hashers, murmur3.New128WithSeed(uint32(capacity)))
	}

	b = bloom{
		Capacity:     capacity,
		ErrorRate:    errorRate,
		hashFuncs:    hashers,
		bitArray:     bitarray.New(baSize),
		bitArraySize: baSize,
		numHashes:    nHashes,
	}

	fmt.Printf("created with array size: %v, using %v hash functions.\n", b.bitArraySize, b.numHashes)

	return b, nil
}

func optimalBitArraySize(n, p float64) int {
	// calculate our bit array size from the probability p and
	m := -(n * math.Log(p)) / (math.Pow(math.Log(2), 2))
	return int(m)
}

func optimalHashFunk(m, n int) int {
	floatM := float64(m)
	floatN := float64(n)

	numFuncs := (floatM / floatN) * math.Log(2)
	return int(numFuncs)
}

func (b bloom) Add(item string) error {
	for i := 0; i < b.numHashes; i++ {
		// Get the hash value
		_, err := b.hashFuncs[i].Write([]byte(item))
		if err != nil {
			return fmt.Errorf("failed to hash the item, %v", err)
		}

		indexPre, _ := b.hashFuncs[i].Sum128()
		// fmt.Printf("indexPre: %v\n", indexPre)
		index := indexPre % uint64(b.bitArraySize)
		// fmt.Printf("index: %v\n", index)

		// Insert to bit array
		_, err = b.bitArray.Put(int(index), 1)
		if err != nil {
			return fmt.Errorf("failed to put bitarray item: %v", index)
		}
		b.hashFuncs[i].Reset()
	}
	return nil
}

// Check answers back whether the item is definitely not in the set (true) or might be in the set (false)
func (b bloom) Check(item string) (bool, error) {
	for i := 0; i < b.numHashes; i++ {
		// Get the hash value
		_, err := b.hashFuncs[i].Write([]byte(item))
		if err != nil {
			return false, fmt.Errorf("failed to hash the item, %v", err)
		}

		indexPre, _ := b.hashFuncs[i].Sum128()
		// fmt.Printf("indexPre: %v\n", indexPre)
		index := indexPre % uint64(b.bitArraySize)
		// fmt.Printf("index: %v\n", index)

		// Check for existence
		res, err := b.bitArray.Get(int(index))
		// fmt.Printf("res: %v\n", res)
		if err != nil {
			return false, fmt.Errorf("failed to check index, %v", err)
		}
		if res == 1 {
			return true, nil
		}
		b.hashFuncs[i].Reset()
	}
	return false, nil
}
