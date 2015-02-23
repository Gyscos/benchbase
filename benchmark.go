package benchbase

import "time"

type Benchmark struct {
	Date   time.Time
	Conf   Configuration
	Result Result
}

func NewBenchmark() *Benchmark {
	return &Benchmark{
		Date:   time.Now(),
		Conf:   NewConfiguration(),
		Result: NewResult(),
	}
}
