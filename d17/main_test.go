package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `.#.
..#
###`

type coord struct{ x, y, z, a int }

func (c coord) neightboors3() []coord {
	n := make([]coord, 0, 26)
	for x := c.x - 1; x <= c.x+1; x++ {
		for y := c.y - 1; y <= c.y+1; y++ {
			for z := c.z - 1; z <= c.z+1; z++ {
				d := coord{x, y, z, 0}
				if d == c {
					continue
				}
				n = append(n, d)
			}
		}
	}
	return n
}

type pocket map[coord]bool

func parseInput(s string) pocket {
	p := make(pocket)
	for y, l := range strings.Split(strings.TrimSpace(s), "\n") {
		for x, b := range l {
			if b == '#' {
				p[coord{x, y, 0, 0}] = true
			}
		}
	}
	return p
}

func (p pocket) activity3() map[coord]int {
	activity := make(map[coord]int)
	for center := range p {
		for _, c := range center.neightboors3() {
			activity[c]++
		}
	}
	return activity
}

func (p pocket) cycle3() pocket {
	next := make(pocket)
	for c, a := range p.activity3() {
		if a == 3 {
			next[c] = true
		} else if a == 2 && p[c] {
			next[c] = true
		}
	}
	return next
}
func (p pocket) cycles3(n int) pocket {
	result := p
	for i := 0; i < n; i++ {
		result = result.cycle3()
	}
	return result
}

func TestFirst(t *testing.T) {
	p := parseInput(example)
	result := p.cycles3(6)
	require.Equal(t, 112, len(result))

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	p = parseInput(string(input))
	result = p.cycles3(6)
	require.Equal(t, 230, len(result))
}

func TestSecond(t *testing.T) {
	p := parseInput(example)
	result := p.cycles4(6)
	require.Equal(t, 848, len(result))

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	p = parseInput(string(input))
	result = p.cycles4(6)
	require.Equal(t, 1600, len(result))
}

func (c coord) neightboors4() []coord {
	n := make([]coord, 0, 80)
	for x := c.x - 1; x <= c.x+1; x++ {
		for y := c.y - 1; y <= c.y+1; y++ {
			for z := c.z - 1; z <= c.z+1; z++ {
				for a := c.a - 1; a <= c.a+1; a++ {
					d := coord{x, y, z, a}
					if d == c {
						continue
					}
					n = append(n, d)
				}
			}
		}
	}
	return n
}

func (p pocket) activity4() map[coord]int {
	activity := make(map[coord]int)
	for center := range p {
		for _, c := range center.neightboors4() {
			activity[c]++
		}
	}
	return activity
}

func (p pocket) cycle4() pocket {
	next := make(pocket)
	for c, a := range p.activity4() {
		if a == 3 {
			next[c] = true
		} else if a == 2 && p[c] {
			next[c] = true
		}
	}
	return next
}
func (p pocket) cycles4(n int) pocket {
	result := p
	for i := 0; i < n; i++ {
		result = result.cycle4()
	}
	return result
}
