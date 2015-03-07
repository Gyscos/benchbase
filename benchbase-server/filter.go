package main

import (
	"encoding/json"
	"regexp"
	"strconv"
)

var patterns = []struct {
	pattern  *regexp.Regexp
	minMatch int
	filter   func([]string, []string) []int
}{
	{
		regexp.MustCompile(`^=([a-zA-Z0-9\.]+)$`),
		2,
		func(values []string, matches []string) []int {
			for i, v := range values {
				if v == matches[1] {
					return []int{i}
				}
			}
			return nil
		},
	},
	{
		regexp.MustCompile(`^!=([a-zA-Z0-9\.]+)$`),
		2,
		func(values []string, matches []string) []int {
			var result []int
			for i, v := range values {
				if v != matches[1] {
					result = append(result, i)
				}
			}
			return invert(result)
		},
	},
	{
		regexp.MustCompile(`^<([a-zA-Z0-9\.]+)$`),
		2,
		func(values []string, matches []string) []int {
			var result []int
			for i, v := range values {
				if Less(v, matches[1]) {
					result = append(result, i)
				} else {
					break
				}
			}
			return invert(result)
		},
	},
	{
		regexp.MustCompile(`^>([a-zA-Z0-9\.]+)$`),
		2,
		func(values []string, matches []string) []int {
			var result []int
			n := len(values)
			for i := n - 1; i >= 0; i-- {
				v := values[i]
				if Less(matches[1], v) {
					result = append(result, i)
				} else {
					break
				}
			}
			// Now invert
			return invert(result)
			// return result
		},
	},
	{
		regexp.MustCompile(`^<=([a-zA-Z0-9\.]+)$`),
		2,
		func(values []string, matches []string) []int {
			var result []int
			for i, v := range values {
				if LessEq(v, matches[1]) {
					result = append(result, i)
				} else {
					break
				}
			}
			return invert(result)
		},
	},
	{
		regexp.MustCompile(`^>=([a-zA-Z0-9\.]+)$`),
		2,
		func(values []string, matches []string) []int {
			var result []int
			n := len(values)
			for i := n - 1; i >= 0; i-- {
				v := values[i]
				if LessEq(matches[1], v) {
					result = append(result, i)
				} else {
					break
				}
			}
			// Now invert
			return invert(result)
			// return result
		},
	},
}

type SpecFilter func([]string) []int

func MakeFilters(description string) (map[string]SpecFilter, error) {
	filters := make(map[string]SpecFilter)

	if description == "" {
		return filters, nil
	}

	var filterData map[string]string
	err := json.Unmarshal([]byte(description), &filterData)
	if err != nil {
		return nil, err
	}

	for k, data := range filterData {
		filters[k] = MakeFilter(data)
	}

	return filters, nil
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

func MakeFilter(data string) SpecFilter {
	for _, p := range patterns {
		matches := p.pattern.FindStringSubmatch(data)
		if len(matches) < p.minMatch {
			continue
		}
		return func(values []string) []int {
			return p.filter(values, matches)
		}
	}

	return func([]string) []int {
		return nil
	}
}
