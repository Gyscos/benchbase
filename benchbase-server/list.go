package main

func intersection(a, b []int) []int {
	var result []int

	m, n := len(a), len(b)
	var i, j int
	for i < m && j < n {
		ai, bj := a[i], b[j]
		if ai < bj {
			i++
		} else if ai > bj {
			j++
		} else {
			result = append(result, ai)
			i++
		}
	}

	return result
}
