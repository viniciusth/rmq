package rmq

import "golang.org/x/exp/constraints"

type RMQNaive[T constraints.Integer | constraints.Float] struct {
  n int
  arr []T
}

func NewRMQNaive[T constraints.Integer | constraints.Float](arr []T) *RMQNaive[T] {
  return &RMQNaive[T]{
    n:   len(arr),
    arr: arr,
  }
}

func (rmq *RMQNaive[T]) Query(l, r int) int {
  if l < 0 || r >= rmq.n || l > r {
    panic("invalid range")
  }

  minIndex := l
  for i := l + 1; i <= r; i++ {
    if rmq.arr[i] < rmq.arr[minIndex] {
      minIndex = i
    }
  }

  return minIndex
}

