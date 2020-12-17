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

type coord struct{ x, y, z int }

func (c coord) neightboors() []coord {
	n := make([]coord, 0, 26)
	for x := c.x - 1; x <= c.x+1; x++ {
		for y := c.y - 1; y <= c.y+1; y++ {
			for z := c.z - 1; z <= c.z+1; z++ {
				d := coord{x, y, z}
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
				p[coord{x, y, 0}] = true
			}
		}
	}
	return p
}

func (p pocket) activity() map[coord]int {
	activity := make(map[coord]int)
	for center := range p {
		for _, c := range center.neightboors() {
			activity[c]++
		}
	}
	return activity
}

func (p pocket) cycle() pocket {
	next := make(pocket)
	for c, a := range p.activity() {
		if a == 3 {
			next[c] = true
		} else if a == 2 && p[c] {
			next[c] = true
		}
	}
	return next
}
func (p pocket) cycles(n int) pocket {
	result := p
	for i := 0; i < n; i++ {
		result = result.cycle()
	}
	return result
}

func TestFirst(t *testing.T) {
	p := parseInput(example)
	result := p.cycles(6)
	require.Equal(t, 112, len(result))

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	p = parseInput(string(input))
	result = p.cycles(6)
	require.Equal(t, 112, len(result))
}
