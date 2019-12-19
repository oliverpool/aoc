package d18

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type coord struct {
	x, y int
}

func (c coord) neighbors() []coord {
	return []coord{
		{c.x, c.y - 1},
		{c.x + 1, c.y},
		{c.x, c.y + 1},
		{c.x - 1, c.y},
	}
}

var notFound = coord{-1, -1}

type rawPixel byte

const (
	entrance = rawPixel('@')
)

func (rp rawPixel) String() string {
	return string(rp)
}

func (rp rawPixel) isKey() bool {
	return rp >= 'a' && rp <= 'z'
}

func (rp rawPixel) Upper() rawPixel {
	return rp + 'A' - 'a'
}

type rawMap []string

func (rm rawMap) get(c coord) rawPixel {
	return rawPixel(rm[c.y][c.x])
}

func (rm rawMap) find(rp rawPixel) coord {
	p := rune(rp)
	for y, line := range rm {
		for x, b := range line {
			if b == p {
				return coord{x, y}
			}
		}
	}
	return notFound
}

func (rm rawMap) pois() map[rawPixel]coord {
	pois := make(map[rawPixel]coord)
	for y, line := range rm {
		for x, b := range line {
			if b == '#' || b == '.' {
				continue
			}
			pois[rawPixel(b)] = coord{x, y}
		}
	}
	return pois
}

func (rm rawMap) neighbors(current coord) map[rawPixel]int {
	neighbors := make(map[rawPixel]int)

	distances := map[coord]int{current: 0}
	next := []coord{current}
	for len(next) > 0 {
		current, next = next[0], next[1:]
		d := distances[current]

		for _, c := range current.neighbors() {
			b := rm.get(c)
			if b == '#' {
				continue
			}
			if _, ok := distances[c]; ok {
				// already visited
				continue
			}
			if b == '.' {
				distances[c] = d + 1
				next = append(next, c)
				continue
			}
			neighbors[b] = d + 1
		}
	}

	return neighbors
}

type poiGraph map[rawPixel]map[rawPixel]int

func parseMap(r io.Reader) (poiGraph, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var rm rawMap
	for scanner.Scan() {
		text := scanner.Text()
		rm = append(rm, text)
	}
	pois := rm.pois()
	distances := make(poiGraph, len(pois))
	for rp, c := range pois {
		distances[rp] = rm.neighbors(c)
	}

	return distances, scanner.Err()
}

func (pg poiGraph) copy() poiGraph {
	cp := make(poiGraph, len(pg))
	for i, m := range pg {
		cm := make(map[rawPixel]int, len(m))
		for j, v := range m {
			cm[j] = v
		}
		cp[i] = cm
	}
	return cp
}

func (pg poiGraph) remove(target rawPixel) {
	neighbors, ok := pg[target]
	if !ok {
		return
	}

	for a, dta := range neighbors {
		for b, dtb := range neighbors {
			if b == a {
				continue
			}
			dab, ok := pg[a][b]
			if !ok || dab > dtb+dta {
				pg[a][b] = dtb + dta
			}
		}
	}
	delete(pg, target)
	for a := range neighbors {
		delete(pg[a], target)
	}
}

func (pg poiGraph) explore(start rawPixel) int {
	cpg := pg.copy()
	cpg.remove(start.Upper()) // open all doors
	neighbors := cpg[start]

	if len(neighbors) == 0 {
		return 0
	}

	cpg.remove(start)
	best := -1
	for n, d := range neighbors {
		if !n.isKey() {
			continue
		}
		current := cpg.explore(n) + d
		if best == -1 || current < best {
			best = current
		}
	}

	return best
}

func (pg poiGraph) exploreBetter(start rawPixel, current, best int) int {
	cpg := pg.copy()
	cpg.remove(start.Upper()) // open all doors
	neighbors := cpg[start]

	if len(neighbors) == 0 {
		return current
	}

	cpg.remove(start)
	// fmt.Println(start, "start search", current, best)
	for n, d := range neighbors {
		if !n.isKey() {
			continue
		}
		if best != -1 && current+d >= best {
			// fmt.Println(start, "skip bigger", n, current+d, best)
			continue
		}
		// fmt.Println(start, "explore", n)
		c := cpg.exploreBetter(n, current+d, best)
		// fmt.Println(start, "explored", n, c)
		if best == -1 || c < best {
			// fmt.Println(start, "better found", n, c, best)
			best = c
		}
	}
	// fmt.Println(start, "best found", best)

	return best
}

func TestParseMap(t *testing.T) {
	cc := []struct {
		input           string
		keyCount        int
		originNeighbors int
		steps           int
	}{
		{
			`
#########
#b.A.@.a#
#########`,
			4,
			2,
			8,
		},
		{
			`
########################
#f.D.E.e.C.b.A.@.a.B.c.#
######################.#
#d.....................#
########################`,
			12,
			2,
			86,
		},
		{
			`
########################
#...............b.C.D.f#
#.######################
#.....@.a.B.c.d.A.e.F.g#
########################`,
			13,
			2,
			132,
		},
		// 		{
		// 			`
		// #################
		// #i.G..c...e..H.p#
		// ########.########
		// #j.A..b...f..D.o#
		// ########@########
		// #k.E..a...g..B.n#
		// ########.########
		// #l.F..d...h..C.m#
		// #################`,
		// 			25,
		// 			8,
		// 			136,
		// 		},
		{
			`
########################
#@..............ac.GI.b#
###d#e#f################
###A#B#C################
###g#h#i################
########################`,
			15,
			4,
			81,
		},
	}
	for _, c := range cc {
		t.Run(c.input, func(t *testing.T) {
			a := assert.New(t)

			m, err := parseMap(strings.NewReader(strings.TrimSpace(c.input)))
			a.NoError(err)
			a.NotNil(m)
			a.Len(m, c.keyCount)
			a.Len(m['@'], c.originNeighbors)

			t.Log(m)

			d := m.exploreBetter('@', 0, -1)
			a.Equal(c.steps, d)
			// fmt.Println("done", c.steps, d)
			// fmt.Println()
		})
	}
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	pg, err := parseMap(f)
	a.NoError(err)

	d := pg.exploreBetter('@', 0, -1)
	a.Equal(2946, d)
}
