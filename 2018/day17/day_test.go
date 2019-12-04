package day_test

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"
)

type coord struct {
	x int
	y int
}

type soil struct {
	isClay         bool
	rightIsBlocked bool
	leftIsBlocked  bool
	parent         *coord
}

type WaterOrigin int

const (
	Unknown WaterOrigin = iota
	Top
	Left
	Right
)

type San

func parseGround(t *testing.T, in io.Reader) (map[coord]*soil, int) {
	sc := bufio.NewScanner(in)

	ground := make(map[coord]*soil)
	max_y := 0

	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "x=") {
			var x, ymin, ymax int
			fmt.Sscanf(line, "x=%d, y=%d..%d", &x, &ymin, &ymax)
			for y := ymin; y <= ymax; y++ {
				ground[coord{x, y}] = &soil{isClay: true}
			}
			if ymax > max_y {
				max_y = ymax
			}
		} else {
			var xmin, xmax, y int
			fmt.Sscanf(line, "y=%d, x=%d..%d", &y, &xmin, &xmax)
			for x := xmin; x <= xmax; x++ {
				ground[coord{x, y}] = &soil{isClay: true}
			}
			if y > max_y {
				max_y = y
			}
		}
	}

	if sc.Err() != nil {
		t.Errorf("could not scan: %v", sc.Err())
	}
	return ground, max_y
}

func TestOne(t *testing.T) {
	example := strings.NewReader(`x=495, y=2..7
y=7, x=495..501
x=501, y=3..7
x=498, y=2..4
x=506, y=1..2
x=498, y=10..13
x=504, y=10..13
y=13, x=498..504`)

	ground, ymax := parseGround(t, example)
	toExplore := []coord{{500, 1}}
	var current coord
	for len(toExplore) > 0 {
		current, toExplore = toExplore[0], toExplore[1:]
		g := ground[current]
		if g == nil {
			g = &soil{}
			ground[current] = g
		}
		if current.y >= ymax {
			continue
		}
		under := coord{current.x, current.y + 1}
		gu := ground[under]
		if !gu.isClay && !gu.leftIsBlocked {
			toExplore = append(toExplore, under)
		} else {
			left := coord{current.x - 1, current.y}
			gl := ground[left]
			if !gl.isClay && gl.waterSettled == nil {
				toExplore = append(toExplore, left)
			} else {
				t.Fatal("left", current)
			}
			right := coord{current.x + 1, current.y}
			gr := ground[right]
			if !gr.isClay && gr.waterSettled == nil {
				toExplore = append(toExplore, right)
			} else {
				t.Fatal("right", current)
			}
		}
		// t.Fatal(g)
	}
	t.Fatal(ground, ymax)
}
