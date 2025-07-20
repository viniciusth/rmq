package rmq

import "golang.org/x/exp/constraints"

// Total space: O(n log n)

type RMQLog[T constraints.Integer | constraints.Float] struct {
	n   int
	log []int
	// O(n log n) memory
	st  [][]int
	arr []T
	// Comparator (a is better than b)
	less func(T, T) bool
}

// Preprocessing: O(n log n) time, O(n log n) space
func NewRMQLog[T constraints.Integer | constraints.Float](arr []T, less func(a, b T) bool) *RMQLog[T] {
	n := len(arr)
	log := make([]int, n+1)
	for i := 2; i <= n; i++ {
		log[i] = log[i/2] + 1
	}

	k := log[n]
	st := make([][]int, k+1)
	for i := range st {
		st[i] = make([]int, n)
	}

	for i := range n {
		st[0][i] = i
	}

	for j := 1; j <= k; j++ {
		for i := 0; i+(1<<j) <= n; i++ {
			idx1 := st[j-1][i]
			idx2 := st[j-1][i+(1<<(j-1))]
			if less(arr[idx1], arr[idx2]) {
				st[j][i] = idx1
			} else {
				st[j][i] = idx2
			}
		}
	}

	return &RMQLog[T]{n: n, log: log, st: st, arr: arr, less: less}
}

// Query time: O(1)
func (rmq *RMQLog[T]) Query(l, r int) int {
	if l < 0 || r >= rmq.n || l > r {
		panic("invalid range")
	}
	j := rmq.log[r-l+1]
	idx1 := rmq.st[j][l]
	idx2 := rmq.st[j][r-(1<<j)+1]
	if rmq.less(rmq.arr[idx1], rmq.arr[idx2]) {
		return idx1
	}
	return idx2
}
