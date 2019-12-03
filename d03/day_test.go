package day

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findLowerLeftCross(wire1, wire2 []string) int {
	w1sx, w1sy := constructSegments(wire1)
	w2sx, w2sy := constructSegments(wire2)
	c1 := findLowestCross(w1sx, w2sy)
	c2 := findLowestCross(w2sx, w1sy)
	if c1 == -1 {
		return c2
	} else if c2 == -1 || c1 < c2 {
		return c1
	} else {
		return c2
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func findLowestCross(sx, sy []segment) int {
	min := -1
	for _, x := range sx {
		for _, y := range sy {
			if y.base < x.start || y.base > x.end {
				continue
			}
			if x.base < y.start || x.base > y.end {
				continue
			}

			// found !
			distance := abs(x.base) + abs(y.base)
			if distance == 0 {
				// ignore start
				continue
			}
			if min == -1 || distance < min {
				min = distance
			}
		}
	}
	return min
}

type segment struct {
	base, start, end int
	distance         int
}

func constructSegments(wire []string) ([]segment, []segment) {
	sx := make([]segment, 0, len(wire)/2)
	sy := make([]segment, 0, len(wire)/2)

	x, y := 0, 0
	var s segment
	var distance int
	for _, w := range wire {
		delta, _ := strconv.Atoi(w[1:])
		distance += delta
		switch w[0] {
		case 'R':
			s, y = segment{x, y, y + delta, distance}, y+delta
			sy = append(sy, s)
		case 'L':
			s, y = segment{x, y - delta, y, -distance}, y-delta
			sy = append(sy, s)
		case 'D':
			s, x = segment{y, x - delta, x, -distance}, x-delta
			sx = append(sx, s)
		case 'U':
			s, x = segment{y, x, x + delta, distance}, x+delta
			sx = append(sx, s)
		default:
			panic("unknown direction " + w)
		}
	}
	return sx, sy
}

func TestLowerLeftCross(t *testing.T) {
	cc := []struct {
		wire1    []string
		wire2    []string
		distance int
	}{
		{
			[]string{"R8", "U5", "L5", "D3"},
			[]string{"U7", "R6", "D4", "L4"},
			6,
		}, {
			[]string{"R75", "D30", "R83", "U83", "L12", "D49", "R71", "U7", "L72"},
			[]string{"U62", "R66", "U55", "R34", "D71", "R55", "D58", "R83"},
			159,
		}, {
			[]string{"R98", "U47", "R26", "D63", "R33", "U87", "L62", "D20", "R33", "U53", "R51"},
			[]string{"U98", "R91", "D20", "R16", "D67", "R40", "U7", "R15", "U6", "R7"},
			135,
		},
	}
	for _, c := range cc {
		t.Run("findLowerLeftCross ", func(t *testing.T) {
			a := assert.New(t)
			distance := findLowerLeftCross(c.wire1, c.wire2)
			a.Equal(c.distance, distance)
		})
	}
}

func SplitByte(sep byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := bytes.IndexByte(data, sep); i >= 0 {
			// We have a full line.
			return i + 1, data[0:i], nil
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}

		// Request more data.
		return 0, nil, nil
	}
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var wires [][]string
	for scanner.Scan() {
		text := scanner.Text()
		wires = append(wires, strings.Split(text, ","))
	}
	a.NoError(scanner.Err())

	a.Equal(557, findLowerLeftCross(wires[0], wires[1]))
}

func TestFewerCombined(t *testing.T) {
	cc := []struct {
		wire1    []string
		wire2    []string
		distance int
	}{
		{
			[]string{"R8", "U5", "L5", "D3"},
			[]string{"U7", "R6", "D4", "L4"},
			30,
		}, {
			[]string{"R75", "D30", "R83", "U83", "L12", "D49", "R71", "U7", "L72"},
			[]string{"U62", "R66", "U55", "R34", "D71", "R55", "D58", "R83"},
			610,
		}, {
			[]string{"R98", "U47", "R26", "D63", "R33", "U87", "L62", "D20", "R33", "U53", "R51"},
			[]string{"U98", "R91", "D20", "R16", "D67", "R40", "U7", "R15", "U6", "R7"},
			410,
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			distance := findFewerCombined(c.wire1, c.wire2)
			a.Equal(c.distance, distance)
		})
	}
}

func findFewerCombined(wire1, wire2 []string) int {
	w1sx, w1sy := constructSegments(wire1)
	w2sx, w2sy := constructSegments(wire2)
	c1 := findLowestCombined(w1sx, w2sy)
	c2 := findLowestCombined(w2sx, w1sy)
	if c1 == -1 {
		return c2
	} else if c2 == -1 || c1 < c2 {
		return c1
	} else {
		return c2
	}
}

func findLowestCombined(sx, sy []segment) int {
	min := -1
	for _, x := range sx {
		for _, y := range sy {
			if y.base < x.start || y.base > x.end {
				continue
			}
			if x.base < y.start || x.base > y.end {
				continue
			}
			if x.base == 0 && y.base == 0 {
				continue
			}

			distance := abs(x.distance)
			if x.distance > 0 {
				// x.end is the end (ie where x.distance is correct)
				distance -= x.end - y.base
			} else {
				// x.end is the start (ie distance is totally wrong)
				distance -= y.base - x.start
			}

			distance += abs(y.distance)
			if y.distance > 0 {
				// y.end is the end (ie where y.distance is correct)
				distance -= y.end - x.base
			} else {
				// y.end is the start (ie distance is totally wrong)
				distance -= x.base - y.start
			}

			if min == -1 || distance < min {
				min = distance
			}
		}
	}
	return min
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var wires [][]string
	for scanner.Scan() {
		text := scanner.Text()
		wires = append(wires, strings.Split(text, ","))
	}
	a.NoError(scanner.Err())

	a.Equal(56410, findFewerCombined(wires[0], wires[1]))
}
