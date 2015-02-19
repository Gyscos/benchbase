package main

type Filter func(BenchMark) bool

var TrueFilter Filter = func(b BenchMark) bool {
	return true
}

func MakeFilter(rule string) Filter {

	return TrueFilter
}
