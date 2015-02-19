package main

// Subject describes the subject of a benchmark.
type Subject struct {
	// SVN revision used in the benchmark
	Rev string
	// Host the benchmark was run on
	Host string
	// Number of threads used in the benchmark
	Threads int
	// Number of threads used in the benchmark divided by number of host CPUs.
	ThreadPerCPU int
}
