package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"github.com/viniciusth/rmq"
)

func readInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	var n int
	for _, c := range scanner.Bytes() {
		n = n*10 + int(c-'0')
	}
	return n
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriterSize(os.Stderr, 4096)
	scanner.Split(bufio.ScanWords)

	n := readInt(scanner)
	arr := make([]int, n)
	for i := range n {
		arr[i] = readInt(scanner)
	}
	rmq := rmq.NewRMQLog(arr)

	m := readInt(scanner)
	for range m {
		l := readInt(scanner)
		r := readInt(scanner)
		if l > r {
			l, r = r, l
		}
		result := rmq.Query(l, r)
		fmt.Fprintln(writer, arr[result])
	}
	writer.Flush()

	PrintMemUsage()
}
