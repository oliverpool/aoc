package day_test

import (
	"testing"
)

// (((x+10)*y+s)*(x+10) /100) % 10 -5

type cell struct {
	x, y int
}

func (c cell) power(serial int) int {
	return (((c.x+10)*c.y+serial)*(c.x+10)/100)%10 - 5
}

func (c cell) cumPower(serial, size int) int {
	sum := 0
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			sum += cell{c.x + x, c.y + y}.power(serial)
		}
	}
	return sum
}

func TestOne(t *testing.T) {
	cc := []struct {
		cell   cell
		serial int
		level  int
	}{
		{cell{122, 79}, 57, -5},
		{cell{217, 196}, 39, 0},
		{cell{101, 153}, 71, 4},
	}

	for _, c := range cc {
		if c.cell.power(c.serial) != c.level {
			t.Errorf("%v with serial %d should have level of %d, got %d", c.cell, c.serial, c.level, c.cell.power(c.serial))
		}
	}

	if m := maxPower(42); m != (cell{21, 61}) {
		t.Errorf("max power of serial 42 should have been at 21,61, got %v", m)
	}

	t.Fatal(maxPower(8868))
}

func maxPower(serial int) cell {
	max := 9 * -10
	var maxCell cell
	for x := 2; x <= 299; x++ {
		for y := 2; y <= 299; y++ {
			current := cell{x - 1, y - 1}.cumPower(serial, 3)
			if current >= max {
				max = current
				maxCell = cell{x - 1, y - 1}
			}
		}
	}
	return maxCell
}

type smartsquare map[cell]int

func NewSmartSquare(size, serial int) smartsquare {
	ss := make(smartsquare, size*size)
	for x := 1; x <= size; x++ {
		for y := 1; y <= size; y++ {
			c := cell{x, y}
			ss[c] = c.power(serial) + ss[cell{x - 1, y}] + ss[cell{x, y - 1}] - ss[cell{x - 1, y - 1}]
		}
	}
	return ss
}

func (ss smartsquare) powerAt(x, y int, size int) int {
	return ss[cell{x - 1, y - 1}] - ss[cell{x - 1, y + size - 1}] + ss[cell{x + size - 1, y + size - 1}] - ss[cell{x + size - 1, y - 1}]
}

func (ss smartsquare) bestCell(size int) (cell, int) {
	max := size * -10
	var maxCell cell
	for x := 1; x <= 300-size+1; x++ {
		for y := 1; y <= 300-size+1; y++ {
			current := ss.powerAt(x, y, size)
			if current >= max {
				max = current
				maxCell = cell{x, y}
			}
		}
	}

	return maxCell, max
}

func (ss smartsquare) bestFlexibleCell() (cell, int) {
	maxSize := 300
	max := maxSize * -5
	var maxCell cell
	var size int
	for s := maxSize; s > 0; s-- {
		if s*s*5 < max {
			break
		}
		c, m := ss.bestCell(s)
		if m > max {
			max = m
			maxCell = c
			size = s
		}
	}
	return maxCell, size
}

func TestTwo(t *testing.T) {
	cc := []struct {
		cell   cell
		serial int
		size   int
		level  int
	}{
		{cell{90, 269}, 18, 16, 113},
		{cell{232, 251}, 42, 12, 119},
	}

	for _, c := range cc {
		ss := NewSmartSquare(300, c.serial)
		if ce, p := ss.bestCell(c.size); p != c.level || ce != c.cell {
			t.Errorf("%v with serial %d should have level of %d, got %d", c.cell, c.serial, c.level, p)
		}
		if ce, p := ss.bestFlexibleCell(); p != c.size || ce != c.cell {
			t.Errorf("%v with serial %d should have level of %d, got %d", c.cell, c.serial, c.level, p)
		}
	}

	ss := NewSmartSquare(300, 8868)
	t.Fatal(ss.bestFlexibleCell())
}
