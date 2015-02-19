package main

// Configuration describes the kind of test used in a benchmark.
// The configuration defines the fields in the result times, so
// two results can only be compared if they have the same configuration.
type Configuration struct {
	// TRUE if all URLs go through the Analyze API. Results are still sorted by type.
	ForceAnalyze bool
	// Depth of the time data to keep. 0 for unlimited.
	Depth int
}
