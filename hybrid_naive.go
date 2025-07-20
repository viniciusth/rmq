package rmq

import (
	"golang.org/x/exp/constraints"
)

// Total space: O(n)

type RMQHybridNaive[T constraints.Integer | constraints.Float] struct {
	n   int
	arr []T
	// O(n / log n) memory (pointers and structs)
	stables []*RMQNaive[T]
	// O(n / log n) memory
	topArr []T
	// O(n) memory
	topST    *RMQLog[T] // This is the RMQLog for the top array
	blockLen int
}

// Preprocessing: O(n) time, O(n) space
func NewRMQHybridNaive[T constraints.Integer | constraints.Float](arr []T) *RMQHybridNaive[T] {
	n := len(arr)

	blockLen := 1
	for (1 << blockLen) < n {
		blockLen++
	}

	blockCount := (n + blockLen - 1) / blockLen
	stables := make([]*RMQNaive[T], blockCount)
	topArr := make([]T, blockCount)
	for i := range blockCount {
		blockStart := i * blockLen
		blockEnd := min((i+1)*blockLen, n)
		stables[i] = NewRMQNaive(arr[blockStart:blockEnd])
		topArr[i] = arr[stables[i].Query(0, blockEnd-blockStart-1)+blockStart]
	}

	// Create the RMQLog for the top array
	rmqTop := NewRMQLog(topArr)

	return &RMQHybridNaive[T]{
		n:        n,
		arr:      arr,
		stables:  stables,
		topArr:   topArr,
		topST:    rmqTop,
		blockLen: blockLen,
	}
}

// Query time: O(log n)
func (rmq *RMQHybridNaive[T]) Query(l, r int) int {
	if l < 0 || r >= rmq.n || l > r {
		panic("invalid range")
	}

	// Find the block indices
	leftBlock := l / rmq.blockLen
	rightBlock := r / rmq.blockLen

	// If both l and r are in the same block, we can use the RMQLog directly
	if leftBlock == rightBlock {
		return rmq.stables[leftBlock].Query(l%rmq.blockLen, r%rmq.blockLen) + leftBlock*rmq.blockLen
	}

	// Otherwise, we need to query the left block, the middle blocks, and the right block
	minIndex := rmq.stables[leftBlock].Query(l%rmq.blockLen, rmq.blockLen-1) + leftBlock*rmq.blockLen
	if leftBlock != rightBlock-1 {
		// If there are middle blocks, we need to query the top RMQLog
		blockIndex := rmq.topST.Query(leftBlock+1, rightBlock-1)
		midIndex := rmq.stables[blockIndex].Query(0, rmq.blockLen-1) + (blockIndex * rmq.blockLen)
		if rmq.arr[midIndex] < rmq.arr[minIndex] {
			minIndex = midIndex
		}
	}
	rightBlockMinIndex := rmq.stables[rightBlock].Query(0, r%rmq.blockLen) + rightBlock*rmq.blockLen
	if rmq.arr[rightBlockMinIndex] < rmq.arr[minIndex] {
		minIndex = rightBlockMinIndex
	}

	return minIndex
}
