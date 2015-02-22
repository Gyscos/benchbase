package benchbase

import "strings"

type Result map[string]float64

func NewResult() Result {
	return make(map[string]float64)
}

func (r Result) Filter(depth int) {
	if depth == 0 {
		return
	}

	for component, time := range r {
		if strings.Count(component, ".") > depth {
			if strings.HasSuffix(component, ".total") {
				r[strings.TrimSuffix(component, ".total")] = time
			}
			delete(r, component)
		}
	}
}
