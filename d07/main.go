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
