package main

// Configuration describes the kind of test used in a benchmark.
// Two results can only be compared if they have the same configuration.
type Configuration struct {
	ForceAnalyze bool
	Depth        int
}
