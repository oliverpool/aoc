package d24

import (
	"bufio"
	"fmt"
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

type coord3d struct {
	x, y, z int
}

func (c coord3d) outer() coord3d {
	if c.x == -1 {
		return coord3d{1, 2, c.z - 1}
	}
	if c.x == 6 {
		return coord3d{3, 2, c.z - 1}
	}
	return c
}

func (c coord3d) neighors() []coord3d {
	n := make([]coord3d, 0, 4)

	// vertical
	switch c.y {
	case 0:
		n = append(n, coord3d{2, 1, c.z - 1})
		n = append(n, coord3d{c.x, c.y + 1, c.z})
	case 1:
		n = append(n, coord3d{c.x, c.y - 1, c.z})
		if c.x == 2 {
			// middle
			n = append(n, coord3d{0, 0, c.z + 1})
			n = append(n, coord3d{1, 0, c.z + 1})
			n = append(n, coord3d{2, 0, c.z + 1})
			n = append(n, coord3d{3, 0, c.z + 1})
			n = append(n, coord3d{4, 0, c.z + 1})
		} else {
			n = append(n, coord3d{c.x, c.y + 1, c.z})
		}
	case 2:
		n = append(n, coord3d{c.x, c.y - 1, c.z})
		n = append(n, coord3d{c.x, c.y + 1, c.z})
	case 3:
		if c.x == 2 {
			// middle
			n = append(n, coord3d{0, 4, c.z + 1})
			n = append(n, coord3d{1, 4, c.z + 1})
			n = append(n, coord3d{2, 4, c.z + 1})
			n = append(n, coord3d{3, 4, c.z + 1})
			n = append(n, coord3d{4, 4, c.z + 1})
		} else {
			n = append(n, coord3d{c.x, c.y - 1, c.z})
		}
		n = append(n, coord3d{c.x, c.y + 1, c.z})
	case 4:
		n = append(n, coord3d{c.x, c.y - 1, c.z})
		n = append(n, coord3d{2, 3, c.z - 1})
	}

	// horizontal
	switch c.x {
	case 0:
		n = append(n, coord3d{1, 2, c.z - 1})
		n = append(n, coord3d{c.x + 1, c.y, c.z})
	case 1:
		n = append(n, coord3d{c.x - 1, c.y, c.z})
		if c.y == 2 {
			// middle
			n = append(n, coord3d{0, 0, c.z + 1})
			n = append(n, coord3d{0, 1, c.z + 1})
			n = append(n, coord3d{0, 2, c.z + 1})
			n = append(n, coord3d{0, 3, c.z + 1})
			n = append(n, coord3d{0, 4, c.z + 1})
		} else {
			n = append(n, coord3d{c.x + 1, c.y, c.z})
		}
	case 2:
		n = append(n, coord3d{c.x - 1, c.y, c.z})
		n = append(n, coord3d{c.x + 1, c.y, c.z})
	case 3:
		if c.y == 2 {
			// middle
			n = append(n, coord3d{4, 0, c.z + 1})
			n = append(n, coord3d{4, 1, c.z + 1})
			n = append(n, coord3d{4, 2, c.z + 1})
			n = append(n, coord3d{4, 3, c.z + 1})
			n = append(n, coord3d{4, 4, c.z + 1})
		} else {
			n = append(n, coord3d{c.x - 1, c.y, c.z})
		}
		n = append(n, coord3d{c.x + 1, c.y, c.z})
	case 4:
		n = append(n, coord3d{c.x - 1, c.y, c.z})
		n = append(n, coord3d{3, 2, c.z - 1})
	}

	return n
}

type area3d map[coord3d]bool

func newArea3d(a area) area3d {
	a3 := make(area3d, len(a))
	for c, b := range a {
		a3[coord3d{c.x, c.y, 0}] = b
	}
	return a3
}

func (a area3d) count() int {
	sum := 0
	for _, b := range a {
		if !b {
			continue
		}
		sum++
	}
	return sum
}
func (a area3d) zMinMax() (int, int) {
	min, max := 0, 0
	for c, b := range a {
		if !b {
			continue
		}
		if c.z > max {
			max = c.z
		}
		if c.z < min {
			min = c.z
		}
	}
	return min, max
}

func (a area3d) String() string {
	s := ""
	zmin, zmax := a.zMinMax()
	for z := zmin; z <= zmax; z++ {
		s += "\nLevel " + fmt.Sprint(z) + ":\n"
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				if a[coord3d{x, y, z}] {
					s += "#"
				} else if x == 2 && y == 2 {
					s += "?"
				} else {
					s += "."
				}
			}
			s += "\n"
		}
	}
	return s
}
func (a area3d) next() area3d {
	next := make(area3d, len(a))
	zmin, zmax := a.zMinMax()
	for z := zmin - 1; z < zmax+2; z++ {
		for x := 0; x < 5; x++ {
			for y := 0; y < 5; y++ {
				if x == 2 && y == 2 {
					continue
				}
				c := coord3d{x, y, z}

				neighors := 0
				for _, n := range c.neighors() {
					if a[n] {
						neighors++
					}
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
	}
	return next
}

func TestSecondExample(t *testing.T) {
	a := assert.New(t)

	m, err := parseArea(strings.NewReader(`....#
#..#.
#..##
..#..
#....`))
	a.NoError(err)

	m3 := newArea3d(m)
	a.Equal(8, m3.count())

	for i := 0; i < 10; i++ {
		m3 = m3.next()
	}
	a.Equal(99, m3.count())
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	m, err := parseArea(f)
	a.NoError(err)

	m3 := newArea3d(m)

	for i := 0; i < 200; i++ {
		m3 = m3.next()
	}
	a.Equal(2009, m3.count())
}
