package main

func validOutermost(needle string, contains map[string][]string) int {
	valid := make(map[string]bool, len(contains[needle]))
	for _, c := range contains[needle] {
		valid[c] = true
	}
	l := 0

	for l != len(valid) {
		l = len(valid)
		for sub := range valid {
			for _, c := range contains[sub] {
				valid[c] = true
			}
		}
	}
	return l
}

func countBags(needle string, bags map[string]map[string]int) int {
	counts := make(map[string]int)

	subCount := func(content map[string]int) (int, bool) {
		n := 0
		for b, i := range content {
			if c, ok := counts[b]; ok {
				n += i * (c + 1)
				continue
			}
			return 0, false
		}
		return n, true
	}
	for counts[needle] == 0 {
		for b, content := range bags {
			if c, ok := subCount(content); ok {
				counts[b] = c
				delete(bags, b)
			}
		}
	}

	return counts[needle]
}
