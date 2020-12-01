package main

func findProduct(goal int, expenses ...int) int {
	seen := make(map[int]bool, len(expenses))
	for _, e := range expenses {
		complement := goal - e
		if seen[complement] {
			return complement * e
		}
		seen[e] = true
	}
	return 0
}

func findProduct3(goal int, expenses ...int) int {
	for i, e := range expenses {
		remaining := goal - e
		found := findProduct(remaining, expenses[i+1:]...)
		if found > 0 {
			return found * e
		}
	}
	return 0
}
