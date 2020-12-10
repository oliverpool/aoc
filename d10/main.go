package main

func countSteps(adapters []int) map[int]int {
	current := 0
	steps := make(map[int]int)
	for _, a := range adapters {
		steps[a-current]++
		current = a
	}

	return steps
}
