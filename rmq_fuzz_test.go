package rmq

import (
	"math/rand"
	"testing"
)

func FuzzRMQ(f *testing.F) {
	f.Fuzz(func(t *testing.T, seed int64, n int, q int) {
		if n < 1 || n > 1000 || q < 1 || q > 1000 {
			t.Skip()
		}

		r := rand.New(rand.NewSource(seed))
		arr := make([]int, n)
		for i := range arr {
			arr[i] = r.Intn(1000000) // Random values
		}

		testComparator := Min[int]
		if r.Intn(2) == 1 {
			testComparator = Max[int]
		}

		naive := NewRMQNaive(arr, testComparator)
		log := NewRMQLog(arr, testComparator)
		hybrid := NewRMQHybrid(arr, testComparator)
		hybridNaive := NewRMQHybridNaive(arr, testComparator)

		algos := []struct {
			name string
			rmq  interface{ Query(int, int) int }
		}{
			{"Log", log},
			{"Hybrid", hybrid},
			{"HybridNaive", hybridNaive},
		}

		for qi := 0; qi < q; qi++ {
			l := r.Intn(n)
			rIdx := r.Intn(n-l) + l // Ensure l <= r < n
			naiveIdx := naive.Query(l, rIdx)
			naiveMin := arr[naiveIdx]

			for _, algo := range algos {
				idx := algo.rmq.Query(l, rIdx)
				if idx < l || idx > rIdx || arr[idx] != naiveMin {
					t.Errorf("%s Query(%d, %d) returned invalid index %d (value %d, expected min %d)", algo.name, l, rIdx, idx, arr[idx], naiveMin)
				}
			}
		}
	})
}
