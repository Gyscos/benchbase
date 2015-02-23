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

var patterns = []struct {
	pattern    *regexp.Regexp
	comparator func(a, b string) bool
}{
	{
		regexp.MustCompile(`^([a-zA-Z0-9\.]+)=([a-zA-Z0-9\.]+)$`),
		func(a, b string) bool {
			return a == b
		},
	},
	{
		regexp.MustCompile(`^([a-zA-Z0-9\.]+)!=([a-zA-Z0-9\.]+)$`),
		func(a, b string) bool {
			return a != b
		},
	},
	{
		regexp.MustCompile(`^([a-zA-Z0-9\.]+)<([a-zA-Z0-9\.]+)$`),
		func(a, b string) bool {
			return Less(a, b)
		},
	},
	{
		regexp.MustCompile(`^([a-zA-Z0-9\.]+)>([a-zA-Z0-9\.]+)$`),
		func(a, b string) bool {
			return Less(b, a)
		},
	},
	{
		regexp.MustCompile(`^([a-zA-Z0-9\.]+)<=([a-zA-Z0-9\.]+)$`),
		func(a, b string) bool {
			return LessEq(a, b)
		},
	},
	{
		regexp.MustCompile(`^([a-zA-Z0-9\.]+)>=([a-zA-Z0-9\.]+)$`),
		func(a, b string) bool {
			return LessEq(b, a)
		},
	},
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

func Less(a, b string) bool {
	intA, intB, err := ParseTwoInts(a, b)
	if err != nil {
		return a < b
	} else {
		return intA < intB
	}
}

func LessEq(a, b string) bool {
	return (a == b) || (Less(a, b))
}

// Builds a Filter from the string description
// Examples:
// - host=c3.xlarge
// - rev>=11158
func MakeSimpleFilter(description string) Filter {
	if description == "" {
		return TrueFilter
	}

	for _, p := range patterns {
		matches := p.pattern.FindStringSubmatch(description)
		if len(matches) == 0 {
			continue
		}
		return func(bench *benchbase.Benchmark) bool {
			return p.comparator(bench.Conf[matches[1]], matches[2])
		}
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
	if value == "" {
		return TrueFilter
	}

	return func(bench *benchbase.Benchmark) bool {
		return bench.Conf[spec] == value
	}
}
