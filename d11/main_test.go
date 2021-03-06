package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `L.LL.LL.LL
LLLLLLL.LL
L.L.L..L..
LLLL.LL.LL
L.LL.LL.LL
L.LLLLL.LL
..L.L.....
LLLLLLLLLL
L.LLLLLL.L
L.LLLLL.LL`

type coord struct {
	x, y int
}

func (c coord) neighboors() []coord {
	return []coord{
		{c.x - 1, c.y - 1},
		{c.x - 1, c.y},
		{c.x - 1, c.y + 1},
		{c.x, c.y - 1},
		{c.x, c.y + 1},
		{c.x + 1, c.y - 1},
		{c.x + 1, c.y},
		{c.x + 1, c.y + 1},
	}
}

func (c coord) sightNeighboors(m map[coord]byte) []coord {
	var sn []coord
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 0 {
				continue
			}
			var c2 coord
			looking := true
			for d := 1; looking; d++ {
				c2 = coord{c.x + i*d, c.y + j*d}
				looking = m[c2] == floor
			}
			//			fmt.Println(c2)
			sn = append(sn, c2)
		}
	}
	return sn
}

const (
	floor  = '.'
	empty  = 'L'
	seated = '#'
)

func parseInput(s string) map[coord]byte {
	m := make(map[coord]byte)
	for y, l := range strings.Split(strings.TrimSpace(s), "\n") {
		for x, b := range l {
			m[coord{x, y}] = byte(b)
		}
	}
	return m
}

func stabilize(m map[coord]byte) map[coord]byte {
	changed := true
	for changed {
		changed = false
		neighboors := make(map[coord]int)
		for c, v := range m {
			if v != seated {
				continue
			}
			for _, n := range c.neighboors() {
				neighboors[n]++
			}
		}
		m2 := make(map[coord]byte)
		for c, s := range m {
			if s == empty && neighboors[c] == 0 {
				m2[c] = seated
				changed = true
			} else if s == seated && neighboors[c] >= 4 {
				m2[c] = empty
				changed = true
			} else {
				m2[c] = s
			}
		}
		m = m2
	}
	return m
}
func countOccupied(m map[coord]byte) int {
	c := 0
	for _, s := range m {
		if s == seated {
			c++
		}
	}
	return c
}

func stabilizeSight(m map[coord]byte) map[coord]byte {
	sightNeighboors := make(map[coord][]coord)
	for c := range m {
		sightNeighboors[c] = c.sightNeighboors(m)
	}

	changed := true
	for changed {
		changed = false
		neighboors := make(map[coord]int)
		for c, v := range m {
			if v != seated {
				continue
			}
			for _, n := range sightNeighboors[c] {
				neighboors[n]++
			}
		}
		m2 := make(map[coord]byte)
		for c, s := range m {
			if s == empty && neighboors[c] == 0 {
				m2[c] = seated
				changed = true
			} else if s == seated && neighboors[c] >= 5 {
				m2[c] = empty
				changed = true
			} else {
				m2[c] = s
			}
		}
		m = m2
	}
	return m
}
func TestFirst(t *testing.T) {
	layout := parseInput(example)
	stab := stabilize(layout)
	c := countOccupied(stab)
	require.Equal(t, 37, c)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	layout = parseInput(string(input))
	stab = stabilize(layout)
	c = countOccupied(stab)
	require.Equal(t, 2481, c)
}

func TestSecond(t *testing.T) {
	layout := parseInput(example)
	stab := stabilizeSight(layout)
	c := countOccupied(stab)
	require.Equal(t, 26, c)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	layout = parseInput(string(input))
	stab = stabilizeSight(layout)
	c = countOccupied(stab)
	require.Equal(t, 2227, c)
}
