package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Gyscos/benchbase"
)

type Filter func(*benchbase.Benchmark) bool

var FalseFilter Filter = func(b *benchbase.Benchmark) bool {
	return false
}

var TrueFilter Filter = func(b *benchbase.Benchmark) bool {
	return true
}

func MakeFilter(description string) Filter {
	var filters []Filter
	rules := strings.Split(description, ";")
	for _, rule := range rules {
		filters = append(filters, MakeSimpleFilter(rule))
	}

	return AndFilter(filters...)
}

func ParseTwoInts(a, b string) (int64, int64, error) {
	intA, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	intB, err := strconv.ParseInt(b, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return intA, intB, nil
}

// Builds a Filter from the string description
// Examples:
// - host=c3.xlarge
// - rev>=11158
func MakeSimpleFilter(description string) Filter {
	patterns := []struct {
		pattern *regexp.Regexp
		loader  func(a, b string) Filter
	}{
		{
			regexp.MustCompile(`^([a-zA-Z0-9\.]+)=([a-zA-Z0-9\.]+)$`),
			func(a, b string) Filter {
				return func(bench *benchbase.Benchmark) bool {
					return bench.Subj[a] == b
				}
			},
		},
		{
			regexp.MustCompile(`^([a-zA-Z0-9\.]+)!=([a-zA-Z0-9\.]+)$`),
			func(a, b string) Filter {
				return func(bench *benchbase.Benchmark) bool {
					return bench.Subj[a] != b
				}
			},
		},
		{
			regexp.MustCompile(`^([a-zA-Z0-9\.]+)<([a-zA-Z0-9\.]+)$`),
			func(a, b string) Filter {
				return func(bench *benchbase.Benchmark) bool {
					intA, intB, err := ParseTwoInts(bench.Subj[a], b)
					if err != nil {
						return bench.Subj[a] < b
					} else {
						return intA < intB
					}
				}
			},
		},
		{
			regexp.MustCompile(`^([a-zA-Z0-9\.]+)>([a-zA-Z0-9\.]+)$`),
			func(a, b string) Filter {
				return func(bench *benchbase.Benchmark) bool {
					intA, intB, err := ParseTwoInts(bench.Subj[a], b)
					if err != nil {
						return bench.Subj[a] > b
					} else {
						return intA > intB
					}
				}
			},
		},
		{
			regexp.MustCompile(`^([a-zA-Z0-9\.]+)<=([a-zA-Z0-9\.]+)$`),
			func(a, b string) Filter {
				return func(bench *benchbase.Benchmark) bool {
					intA, intB, err := ParseTwoInts(bench.Subj[a], b)
					if err != nil {
						return bench.Subj[a] <= b
					} else {
						return intA <= intB
					}
				}
			},
		},
		{
			regexp.MustCompile(`^([a-zA-Z0-9\.]+)>=([a-zA-Z0-9\.]+)$`),
			func(a, b string) Filter {
				return func(bench *benchbase.Benchmark) bool {
					intA, intB, err := ParseTwoInts(bench.Subj[a], b)
					if err != nil {
						return bench.Subj[a] >= b
					} else {
						return intA >= intB
					}
				}
			},
		},
	}

	for _, p := range patterns {
		matches := p.pattern.FindStringSubmatch(description)
		if len(matches) == 0 {
			continue
		}
		return p.loader(matches[1], matches[2])
	}
	log.Println("Could not read description")
	return FalseFilter
}

func AndFilter(filters ...Filter) Filter {

	return func(b *benchbase.Benchmark) bool {
		for _, f := range filters {
			if !f(b) {
				return false
			}
		}
		return true
	}
}

// Builds a special kind of filter
func MakeSpecFilter(spec string, value string) Filter {
	return TrueFilter
}
