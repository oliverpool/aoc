package main

func treeEncounters(forest []string, right, bottom int) int {
	x, y := right, bottom
	l := len(forest[0])
	encounters := 0
	for y < len(forest) {
		if forest[y][x] == '#' {
			encounters++
		}
		x = (x + right) % l
		y += bottom
	}
	return encounters
}

func treeEncountersProduct(forest []string) int {
	return treeEncounters(forest, 1, 1) * treeEncounters(forest, 3, 1) * treeEncounters(forest, 5, 1) * treeEncounters(forest, 7, 1) * treeEncounters(forest, 1, 2)
}
