package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/viniciusth/rmq"
)

type RMQ interface {
	Query(l, r int) int
}

type algo struct {
	name string
	new  func([]int) RMQ
}

var algos = map[string]algo{
	"log":          {name: "log", new: func(a []int) RMQ { return rmq.NewRMQLog(a, rmq.Min) }},
	"hybrid_log":   {name: "hybrid_log", new: func(a []int) RMQ { return rmq.NewRMQHybrid(a, rmq.Min) }},
	"hybrid_naive": {name: "hybrid_naive", new: func(a []int) RMQ { return rmq.NewRMQHybridNaive(a, rmq.Min) }},
}

type memMonitor struct {
	maxAlloc uint64
	stop     chan struct{}
}

func newMemMonitor() *memMonitor {
	mm := &memMonitor{stop: make(chan struct{})}
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			if m.Alloc > mm.maxAlloc {
				mm.maxAlloc = m.Alloc
			}
			select {
			case <-mm.stop:
				return
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
	return mm
}

func (mm *memMonitor) Stop() uint64 {
	close(mm.stop)
	return mm.maxAlloc
}

func getCurrentAlloc() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func measureConstruct(newFunc func([]int) RMQ, arr []int) (time.Duration, uint64, uint64, RMQ) {
	runtime.GC()
	mm := newMemMonitor()
	start := time.Now()
	rmqq := newFunc(arr)
	dur := time.Since(start)
	peak := mm.Stop()
	runtime.GC()
	alloc := getCurrentAlloc()
	return dur, peak, alloc, rmqq
}

func measureQuery(rmqq RMQ, queries [][2]int) (time.Duration, uint64, uint64) {
	runtime.GC()
	mm := newMemMonitor()
	start := time.Now()
	for _, q := range queries {
		_ = rmqq.Query(q[0], q[1])
	}
	dur := time.Since(start)
	peak := mm.Stop()
	runtime.GC()
	alloc := getCurrentAlloc()
	return dur, peak, alloc
}

func runBenchmark(algo algo, N, Q, runs int) {
	for run := 0; run < runs; run++ {
		r := rand.New(rand.NewSource(int64(run)))
		arr := make([]int, N)
		for i := range arr {
			arr[i] = r.Int()
		}
		queries := make([][2]int, Q)
		for i := range queries {
			l := r.Intn(N)
			rr := r.Intn(N)
			if l > rr {
				l, rr = rr, l
			}
			queries[i] = [2]int{l, rr}
		}
		ct, cp, ca, rmqq := measureConstruct(algo.new, arr)
		qt, qp, qa := measureQuery(rmqq, queries)
		fmt.Printf("%s,%d,%d,%.0f,%d,%d,%.0f,%d,%d\n",
			algo.name, N, Q,
			float64(ct.Nanoseconds()), cp, ca,
			float64(qt.Nanoseconds()), qp, qa)
	}
}

func main() {
	algoName := flag.String("algo", "", "Algorithm to benchmark")
	n := flag.Int("n", 0, "Number of elements N")
	q := flag.Int("q", 0, "Number of queries Q")
	runs := flag.Int("runs", 3, "Number of runs for averaging")
	flag.Parse()

	if *algoName == "" || *n <= 0 || *q <= 0 {
		fmt.Println("Usage: go run main.go -algo=<algo> -n=<N> -q=<Q> [-runs=<runs>]")
		fmt.Println("Available algos:", algos)
		os.Exit(1)
	}

	algo, ok := algos[*algoName]
	if !ok {
		fmt.Println("Invalid algo:", *algoName)
		os.Exit(1)
	}

	runBenchmark(algo, *n, *q, *runs)
}
