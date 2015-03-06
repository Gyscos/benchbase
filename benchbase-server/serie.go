package main

func consecutive(start, stop int) []int {
	n := stop - start
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = start + i
	}
	return result
}

func invert(list []int) []int {
	n := len(list)
	result := make([]int, n)

	for i := 0; i < n; i++ {
		result[n-i-1] = list[i]
	}

	return result
}
