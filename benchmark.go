package main

import "time"

type BenchMark struct {
	Date    time.Time
	Subj    Subject
	Conf    Configuration
	Results map[string]float64
}
