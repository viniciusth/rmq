all: plot

bench: bench.csv

bench.csv:
	@echo "algo,N,Q,construct_time_ns,construct_peak_bytes,construct_alloc_bytes,query_time_ns,query_peak_bytes,query_alloc_bytes" > bench.csv
	@for algo in log hybrid_log hybrid_naive; do \
		n=100000; \
		while [ $$n -le 100000000 ]; do \
			echo "Benchmarking $$algo with N=$$n Q=1000000"; \
			go run cmd/bench/main.go -algo=$$algo -n=$$n -q=1000000 -runs=5 >> bench.csv; \
			n=$$((n * 12 / 10)); \
		done; \
	done
	@for algo in log hybrid_log hybrid_naive; do \
		q=100000; \
		while [ $$q -le 10000000 ]; do \
			echo "Benchmarking $$algo with N=1000000 Q=$$q"; \
			go run cmd/bench/main.go -algo=$$algo -n=1000000 -q=$$q -runs=5 >> bench.csv; \
			q=$$((q * 12 / 10)); \
		done; \
	done

plot: bench.csv
	python3 plot.py bench.csv

clean:
	rm -f bench.csv
	rm -rf plots 