package day_test

import (
	"fmt"
	"strings"
	"testing"
)

func reduce(polymer []rune) []rune {
	changes := 1
	for changes > 0 {
		changes = 0
		for index := 0; index < len(polymer)-1; index++ {
			if polymer[index] == -polymer[index+1] {
				// t.Log(polymer, index)
				polymer = append(polymer[:index], polymer[index+2:]...)
				// t.Log(polymer)
				changes += 1
			}
		}
	}
	return polymer

}

type point struct {
	index    int
	distance int
}

type coord struct {
	x, y int
}

const infinity = 100000

func TestOne(t *testing.T) {
	in := input
	li := strings.Split(in, "\n")
	coordinates := make([]coord, len(li))

	maxX, maxY := 0, 0
	for i, s := range li {
		fmt.Sscanf(s, "%d, %d", &coordinates[i].y, &coordinates[i].x)
		if coordinates[i].x > maxX {
			maxX = coordinates[i].x
		}
		if coordinates[i].y > maxY {
			maxY = coordinates[i].y
		}
	}

	grid := make([][]*point, maxX+2)

	printGrid := func() {
		return
		for _, line := range grid {
			for _, p := range line {
				if p == nil {
					fmt.Print("*")
				} else if p.index == -1 {
					fmt.Print(".")
				} else {
					fmt.Print(string(rune(p.index) + 'a'))
				}
			}
			fmt.Println()

		}
	}
	for i := range grid {
		grid[i] = make([]*point, maxY+2)
	}
	for i, c := range coordinates {
		grid[c.x][c.y] = &point{i, 1}
	}
	printGrid()
	fmt.Println("before")
	findClosest := func(i, j int) *point {
		top := grid[i-1][j]
		bot := grid[i+1][j]
		lef := grid[i][j-1]
		rig := grid[i][j+1]
		if top == nil && bot == nil && lef == nil && rig == nil {
			return nil
		}
		p := &point{-1, infinity}
		if top != nil {
			p.index = top.index
			p.distance = top.distance + 1
		}
		if bot != nil && bot.distance+1 <= p.distance {
			if bot.distance+1 < p.distance {
				p.index = bot.index
				p.distance = bot.distance + 1
			} else if bot.index != p.index {
				p.index = -1
			}
		}
		if lef != nil && lef.distance+1 <= p.distance {
			if lef.distance+1 < p.distance {
				p.index = lef.index
				p.distance = lef.distance + 1
			} else if lef.index != p.index {
				p.index = -1
			}
		}
		if rig != nil && rig.distance+1 <= p.distance {
			if rig.distance+1 < p.distance {
				p.index = rig.index
				p.distance = rig.distance + 1
			} else if rig.index != p.index {
				p.index = -1
			}
		}
		return p
	}
	found := true
	for found {
		found = false
		for i := 1; i <= maxX; i++ {
			for j := 1; j <= maxY; j++ {

				if c := findClosest(i, j); c != nil {
					if grid[i][j] != nil && grid[i][j].distance <= c.distance {
						if grid[i][j].distance == c.distance && grid[i][j].index != c.index {
							grid[i][j] = &point{-1, c.distance}
							found = true
						}
						continue
					}
					grid[i][j] = c
					found = true
				}
			}
		}
	}
	printGrid()
	surfaces := make(map[int]int, len(coordinates))
	for i := 1; i <= maxX; i++ {
		for j := 1; j <= maxY; j++ {
			p := grid[i][j]
			if i == 1 || i == maxX || j == 1 || j == maxY {
				surfaces[p.index] = infinity
				continue
			}
			surfaces[p.index] += 1
		}
	}
	maxSurface := 0
	for _, s := range surfaces {
		if s >= infinity || s < maxSurface {
			continue
		}
		maxSurface = s
	}
	t.Fatal(maxSurface)

	polymer := make([]rune, len(in))
	for i, p := range in {
		if p >= 'a' {
			polymer[i] = 'a' - p - 1 // negative
		} else {
			polymer[i] = p - 'A' + 1 // positive
		}
	}
	polymer = reduce(polymer)

	minMinification := 11111
	for index := rune(1); index <= rune(27); index++ {
		var rpoly []rune
		for _, p := range polymer {
			if p == index || p == -index {
				continue
			}
			rpoly = append(rpoly, p)
		}
		minif := len(reduce(rpoly))
		if minif < minMinification {
			minMinification = minif
		}
	}
	t.Fatal(minMinification)
}

func TestTwo(t *testing.T) {
	in := input
	maxDistance := 10000
	li := strings.Split(in, "\n")
	coordinates := make([]coord, len(li))

	maxX, maxY := 0, 0
	for i, s := range li {
		fmt.Sscanf(s, "%d, %d", &coordinates[i].y, &coordinates[i].x)
		if coordinates[i].x > maxX {
			maxX = coordinates[i].x
		}
		if coordinates[i].y > maxY {
			maxY = coordinates[i].y
		}
	}

	grid := make([][]*point, maxX+2)

	printGrid := func() {
		for _, line := range grid {
			for _, p := range line {
				if p == nil {
					fmt.Print("*")
				} else if p.index == -1 {
					fmt.Print(".")
				} else if p.distance < maxDistance {
					fmt.Print("#")
				} else {
					fmt.Print("_")
				}
			}
			fmt.Println()
		}
	}
	for i := range grid {
		grid[i] = make([]*point, maxY+2)
	}
	for i, line := range grid {
		for j := range line {
			p := &point{}
			for _, c := range coordinates {
				p.distance += abs(c.x-i) + abs(c.y-j)
			}
			grid[i][j] = p
		}
	}
	size := 0
	for _, line := range grid {
		for _, p := range line {
			if p.distance < maxDistance {
				size += 1
			}
		}
	}
	t.Fatal(size)
	printGrid()
	fmt.Println("before")
}
func abs(v int) int {
	if v > 0 {
		return v
	}
	return -v
}

var inputTest = `1, 1
1, 6
8, 3
3, 4
5, 5
8, 9`

var input = `59, 110
127, 249
42, 290
90, 326
108, 60
98, 168
358, 207
114, 146
242, 170
281, 43
233, 295
213, 113
260, 334
287, 260
283, 227
328, 235
96, 259
232, 177
198, 216
52, 115
95, 258
173, 191
156, 167
179, 135
235, 235
164, 199
248, 180
165, 273
160, 297
102, 96
346, 249
176, 263
140, 101
324, 254
72, 211
126, 337
356, 272
342, 65
171, 295
93, 192
47, 200
329, 239
60, 282
246, 185
225, 324
114, 329
134, 167
212, 104
338, 332
293, 94`
