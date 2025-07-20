package rmq

import "golang.org/x/exp/constraints"

// Total space: O(n)

type RMQNaive[T constraints.Integer | constraints.Float] struct {
	n   int
	arr []T
	// Comparator (a is better than b)
	less func(T, T) bool
}

func NewRMQNaive[T constraints.Integer | constraints.Float](arr []T, less func(a, b T) bool) *RMQNaive[T] {
	return &RMQNaive[T]{
		n:    len(arr),
		arr:  arr,
		less: less,
	}
}

// Query time: O(r - l + 1)
func (rmq *RMQNaive[T]) Query(l, r int) int {
	if l < 0 || r >= rmq.n || l > r {
		panic("invalid range")
	}

	minIndex := l
	for i := l + 1; i <= r; i++ {
		if rmq.less(rmq.arr[i], rmq.arr[minIndex]) {
			minIndex = i
		}
	}

	return minIndex
}
