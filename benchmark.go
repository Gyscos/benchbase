package benchbase

import "time"

type Benchmark struct {
	Date   time.Time
	Subj   Subject
	Conf   Configuration
	Result Result
}

func NewBenchmark() *Benchmark {
	return &Benchmark{
		Date:   time.Now(),
		Subj:   NewSubject(),
		Conf:   NewConfiguration(),
		Result: NewResult(),
	}
}
