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
