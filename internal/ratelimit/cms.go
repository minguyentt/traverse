package ratelimit

import (
	"encoding/binary"
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"traverse/internal/assert"
)

/*
buckets: # of columns or cells in the hash function
depth: # of hash functions
count: 2d matrix to store the incremental count

bitwise for 64bit
1<<64-1

	hash[i] = (h + j * i) % buckets
*/
type countMinSketch struct {
	buckets  uint
	depth    uint
	hashFunc hash.Hash
	counter  [][]uint64
}

func NewCMS(buckets uint, depth uint) (*countMinSketch, error) {
	if buckets <= 0 || depth <= 0 {
		return nil, fmt.Errorf("depth and buckets should be greater than 0")
	}

	cms := &countMinSketch{
		buckets:  buckets,
		depth:    depth,
		hashFunc: fnv.New64(), // lower risk of collisions w/ fnv hash 64bit
	}

	cms.counter = make([][]uint64, depth)

	for i := uint(0); i < depth; i++ {
		cms.counter[i] = make([]uint64, buckets)
	}

	return cms, nil
}

func (c *countMinSketch) coefficients(key []byte) (uint32, uint32) {
	c.hashFunc.Reset()
	_, _ = c.hashFunc.Write(key)

	// sum = 8 bytes long for FNV-64
	hashedBytes := c.hashFunc.Sum(nil)
	assert.Assert(len(hashedBytes) <= 8, "sum of hash output too short: expected at least 8 bytes", "output", hashedBytes, "leng", len(hashedBytes))

	// split into two 32-bit vals for the coefficients
	h := binary.BigEndian.Uint32(hashedBytes[4:8])
	j := binary.BigEndian.Uint32(hashedBytes[0:4])

	return h, j
}

func (c *countMinSketch) bucketPositions(key []byte) []uint {
	pos := make([]uint, c.depth)
	h, j := c.coefficients(key)

	upper := uint(h)
	lower := uint(j)

	for i := uint(0); i < c.depth; i++ {
		pos[i] = (upper + lower*i) % c.buckets
	}

	return pos
}

func (c *countMinSketch) Update(key []byte, count uint64) {
	for i, j := range c.bucketPositions(key) {
		c.counter[i][j] += count
	}
}

// estimate min. frequency for key
func (c *countMinSketch) Estimate(key []byte) uint64 {
	pos := c.bucketPositions(key)
	min := uint64(math.MaxUint64)

	for row, col := range pos {
		val := c.counter[row][col]
		if val < min {
			min = val
		}
	}

	return min
}

func (c *countMinSketch) MergeStreams(cms *countMinSketch) error {
	assert.Assert(c.depth != cms.depth, "merging streams must match matrix depth")
	assert.Assert(c.buckets != cms.buckets, "merging streams must match matrix buckets")

	for i := uint(0); i < c.depth; i++ {
		for j := uint(0); j < c.buckets; j++ {
			c.counter[i][j] += cms.counter[i][j]
		}
	}

	return nil
}
