package d10

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	if x < 0 {
		return -x
	}
	return x
}

func TestGCD(t *testing.T) {
	testCases := []struct {
		a, b, out int
	}{
		{1, 1, 1},
		{1, -1, 1},
		{-1, -1, 1},
		{-1, 1, 1},
		{0, 1, 1},
		{0, 10, 10},
		{5, 10, 5},
	}
	for _, tC := range testCases {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			a.Equal(tC.out, gcd(tC.a, tC.b))
		})
	}
}

type asteroid struct {
	x, y int
}

type direction struct {
	x, y int
}

func (a asteroid) direction(b asteroid) direction {
	x := b.x - a.x
	y := b.y - a.y
	g := gcd(x, y)
	return direction{x / g, y / g}
}

func parseMap(r io.Reader) []asteroid {
	scanner := bufio.NewScanner(r)
	as := make([]asteroid, 0)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, c := range line {
			if c != '#' {
				continue
			}
			as = append(as, asteroid{x, y})
		}
		y++
	}
	return as
}

func TestParseMap(t *testing.T) {
	testCases := []struct {
		in  string
		out []asteroid
	}{
		{in: `.#..#
.....
#####
....#
...##`,
			out: []asteroid{
				{1, 0},
				{4, 0},
				{0, 2},
				{1, 2},
				{2, 2},
				{3, 2},
				{4, 2},
				{4, 3},
				{3, 4},
				{4, 4},
			}},
	}
	for _, tC := range testCases {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			a.Equal(tC.out, parseMap(strings.NewReader(tC.in)))
		})
	}
}

func maxDirections(as []asteroid) int {
	max := -1
	for _, a := range as {
		d := make(map[direction]int)
		for _, b := range as {
			if a == b {
				continue
			}
			d[a.direction(b)]++
		}
		if max == -1 || len(d) > max {
			max = len(d)
		}
	}
	return max
}

func TestMaxDirection(t *testing.T) {
	testCases := []struct {
		in  string
		out int
	}{
		{
			in: `.#..#
.....
#####
....#
...##`,
			out: 8,
		},
		{
			in: `......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####`,
			out: 33,
		},
		{
			in: `

#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.
`,
			out: 35,
		},
		{
			in: `.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..`,
			out: 41,
		},
		{
			in: `.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`,
			out: 210,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.in, func(t *testing.T) {
			a := assert.New(t)
			a.Equal(tC.out, maxDirections(parseMap(strings.NewReader(tC.in))))
		})
	}
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	as := parseMap(f)

	m := maxDirections(as)
	a.Equal(276, m)
}
