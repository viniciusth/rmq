package rmq

import "golang.org/x/exp/constraints"

// Min returns true if a is less than b, for finding minima.
func Min[T constraints.Ordered](a, b T) bool {
	return a < b
}

// Max returns true if a is greater than b, for finding maxima (treats larger as 'better').
func Max[T constraints.Ordered](a, b T) bool {
	return a > b
}
