package d10

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
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

func (d direction) angle() float64 {
	if d.y < 0 && d.x < 0 {
		return math.Atan2(float64(d.y), float64(d.x)) + 5*math.Pi/2
	}
	return math.Atan2(float64(d.y), float64(d.x)) + math.Pi/2
}

func TestAngles(t *testing.T) {
	testCases := []struct {
		d     direction
		angle float64
	}{
		{direction{0, -3}, 0},
		{direction{2, 0}, math.Pi / 2},
		{direction{0, 3}, math.Pi},
		{direction{-3, 0}, 3 * math.Pi / 2},
		{direction{1, -1}, math.Pi / 4},
		{direction{1, 1}, 3 * math.Pi / 4},
		{direction{-1, 1}, 5 * math.Pi / 4},
		{direction{-1, -1}, 7 * math.Pi / 4},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprint(tC.d), func(t *testing.T) {
			a := assert.New(t)
			a.Equal(tC.angle, tC.d.angle())
		})
	}
}

func (a asteroid) direction(b asteroid) direction {
	x := b.x - a.x
	y := b.y - a.y
	g := gcd(x, y)
	return direction{x / g, y / g}
}

func (a asteroid) distance(b asteroid) int {
	x := b.x - a.x
	y := b.y - a.y
	return gcd(x, y)
}

type ByDirection struct {
	as  [][]asteroid
	src asteroid
}

func (a ByDirection) Len() int      { return len(a.as) }
func (a ByDirection) Swap(i, j int) { a.as[i], a.as[j] = a.as[j], a.as[i] }
func (a ByDirection) Less(i, j int) bool {
	return a.src.direction(a.as[i][0]).angle() < a.src.direction(a.as[j][0]).angle()
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

func maxDirections(as []asteroid) (int, asteroid) {
	max := -1
	var maxA asteroid
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
			maxA = a
		}
	}
	return max, maxA
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
			m, _ := maxDirections(parseMap(strings.NewReader(tC.in)))
			a.Equal(tC.out, m)
		})
	}
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	as := parseMap(f)

	m, mA := maxDirections(as)
	a.Equal(276, m)
	a.Equal(asteroid{17, 22}, mA)
}

func vaporizeFrom(a asteroid, as []asteroid, n int) asteroid {
	directions := make(map[direction][]asteroid)
	for _, b := range as {
		if a == b {
			continue
		}
		dir := a.direction(b)
		line := directions[dir]
		dist := a.distance(b)
		var i int
		for i = 0; i < len(line) && dist > a.distance(line[i]); i++ {
		}
		newline := make([]asteroid, 0, len(line)+1)
		newline = append(newline, line[0:i]...)
		newline = append(newline, b)
		newline = append(newline, line[i:]...)
		directions[dir] = newline
	}

	adirections := make([][]asteroid, 0, len(directions))
	for _, line := range directions {
		adirections = append(adirections, line)
	}

	bd := ByDirection{
		as:  adirections,
		src: a,
	}
	sort.Sort(bd)

	adirections = bd.as

	i := 0
	for n > 1 {
		if len(adirections[i]) <= 1 {
			copy(adirections[i:], adirections[i+1:])
			adirections = adirections[:len(adirections)-1]
		} else {
			adirections[i] = adirections[i][1:]
			i++
		}
		n--
		i = i % len(adirections)
	}
	return adirections[i][0]
}

func TestVaporize(t *testing.T) {
	testCases := []struct {
		in  string
		out asteroid
	}{
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
			out: asteroid{8, 2},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.in, func(t *testing.T) {
			a := assert.New(t)
			last := vaporizeFrom(asteroid{11, 13}, parseMap(strings.NewReader(tC.in)), 200)
			a.Equal(tC.out, last)
		})
	}
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	as := parseMap(f)

	last := vaporizeFrom(asteroid{17, 22}, as, 200)
	a.Equal(asteroid{13, 21}, last)
}
