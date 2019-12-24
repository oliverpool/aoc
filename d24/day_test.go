package d24

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

type area map[coord]bool

func parseArea(r io.Reader) (area, error) {
	m := make(area)

	scanner := bufio.NewScanner(r)
	y := 0
	for scanner.Scan() {
		text := scanner.Text()
		for x, b := range text {
			if b == '#' {
				m[coord{x, y}] = true
			}
		}
		y++
	}
	return m, scanner.Err()
}

func (a area) next() area {
	next := make(area, len(a))
	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			c := coord{x, y}
			neighors := 0
			if a[coord{x - 1, y}] {
				neighors++
			}
			if a[coord{x, y - 1}] {
				neighors++
			}
			if a[coord{x + 1, y}] {
				neighors++
			}
			if a[coord{x, y + 1}] {
				neighors++
			}
			if a[c] {
				// die or not
				next[c] = neighors == 1
			} else {
				// check infested or not
				next[c] = neighors == 1 || neighors == 2
			}
		}
	}
	return next
}

func (a area) rating() int {
	sum := 0
	for c, b := range a {
		if !b {
			continue
		}
		sum += 1 << (c.x + c.y*5)
	}
	return sum
}

func TestParse(t *testing.T) {
	cc := []struct {
		input       string
		len, rating int
	}{
		{
			`....#
#..#.
#..##
..#..
#....`, 8, 1205552,
		},
		{
			`.....
.....
.....
#....
.#...`, 2, 2129920,
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			m, err := parseArea(strings.NewReader(c.input))
			a.NoError(err)
			a.Len(m, c.len)
			a.Equal(c.rating, m.rating())
		})
	}
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	m, err := parseArea(f)
	a.NoError(err)

	found := make(map[int]bool)
	for {
		r := m.rating()
		if found[r] {
			a.Equal(32526865, r)
			break
		}
		found[r] = true
		m = m.next()
	}
}
