package rmq

import (
	"golang.org/x/exp/constraints"
)

// Total space: O(n log log n)

type RMQHybridLog[T constraints.Integer | constraints.Float] struct {
	n   int
	arr []T
	// O(n/logn * (logn log log n) = n log log n) memory (across all block sparse tables)
	stables []*RMQLog[T] // we could optimize the log array to be shared across all blocks, but it wouldn't save that much
	// O(n / log n) memory
	topArr []T
	// O(n) memory
	topST    *RMQLog[T]
	blockLen int
	// Comparator (a is better than b)
	less func(T, T) bool
}

// Preprocessing: O(n log log n) time, O(n log log n) space
func NewRMQHybrid[T constraints.Integer | constraints.Float](arr []T, less func(a, b T) bool) *RMQHybridLog[T] {
	n := len(arr)

	blockLen := 1
	for (1 << blockLen) < n {
		blockLen++
	}

	blockCount := (n + blockLen - 1) / blockLen
	stables := make([]*RMQLog[T], blockCount)
	topArr := make([]T, blockCount)
	for i := range blockCount {
		blockStart := i * blockLen
		blockEnd := min((i+1)*blockLen, n)
		stables[i] = NewRMQLog(arr[blockStart:blockEnd], less)
		topArr[i] = arr[stables[i].Query(0, blockEnd-blockStart-1)+blockStart]
	}

	rmqTop := NewRMQLog(topArr, less)

	return &RMQHybridLog[T]{
		n:        n,
		arr:      arr,
		stables:  stables,
		topArr:   topArr,
		topST:    rmqTop,
		blockLen: blockLen,
		less:     less,
	}
}

// Query time: O(1)
func (rmq *RMQHybridLog[T]) Query(l, r int) int {
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
		if rmq.less(rmq.arr[midIndex], rmq.arr[minIndex]) {
			minIndex = midIndex
		}
	}
	rightBlockMinIndex := rmq.stables[rightBlock].Query(0, r%rmq.blockLen) + rightBlock*rmq.blockLen
	if rmq.less(rmq.arr[rightBlockMinIndex], rmq.arr[minIndex]) {
		minIndex = rightBlockMinIndex
	}

	return minIndex
}
