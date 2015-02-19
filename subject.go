package main

// Subject describes the subject of a benchmark.
type Subject struct {
	Rev          string
	Host         string
	Threads      int
	ThreadPerCPU int
}
