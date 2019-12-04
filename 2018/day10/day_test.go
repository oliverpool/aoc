package day_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

type point struct {
	x, y, vx, vy int
}

func (p point) At(t int) point {
	return point{
		x: p.x + t*p.vx,
		y: p.y + t*p.vy,
	}
}

func (p point) Xat(t int) int {
	return p.x + t*p.vx
}

func (p point) Yat(t int) int {
	return p.y + t*p.vy
}

func TestOne(t *testing.T) {
	f, err := os.Open("./input")
	if err != nil {
		t.Fatal(err)
	}

	var pp []point
	s := bufio.NewScanner(f)
	for s.Scan() {
		p := &point{}
		fmt.Sscanf(s.Text(), "position=<%d,%d> velocity=<%d,%d>", &p.x, &p.y, &p.vx, &p.vy)
		pp = append(pp, *p)
	}
	if s.Err() != nil {
		t.Fatal(s.Err())
	}

	delta := 1
	tt := 0
	previous, _ := coherence(pp, tt)
	for delta >= 0 {
		delta = 0
		tt += 1
		current, _ := coherence(pp, tt)
		delta = previous - current
		previous = current
		// t.Log(tt, cx, cy)
	}

	// cx, cy := coherence(pp, tt)
	for t := tt - 5; t < tt+1; t++ {
		printSky(pp, t)
		fmt.Println()
	}
	t.Error(tt - 1)
	// printSky(pp, tt)
}

func coherence(pp []point, t int) (int, int) {
	minx, maxx, miny, maxy := pp[0].Xat(t), pp[0].Xat(t), pp[0].Yat(t), pp[0].Yat(t)

	for _, p := range pp {
		x := p.Xat(t)
		y := p.Yat(t)
		if x < minx {
			minx = x
		}
		if x > maxx {
			maxx = x
		}
		if y < miny {
			miny = y
		}
		if y > maxy {
			maxy = y
		}
	}
	return (maxx - minx), (maxy - miny)
}

func printSky(pp []point, t int) {
	sky := make(map[point]bool)
	minx, maxx, miny, maxy := pp[0].Xat(t), pp[0].Xat(t), pp[0].Yat(t), pp[0].Yat(t)

	for _, p := range pp {
		x := p.Xat(t)
		y := p.Yat(t)
		if x < minx {
			minx = x
		}
		if x > maxx {
			maxx = x
		}
		if y < miny {
			miny = y
		}
		if y > maxy {
			maxy = y
		}
		sky[p.At(t)] = true
	}
	for y := miny - 1; y <= maxy+1; y++ {
		for x := minx - 1; x <= maxx+1; x++ {
			if sky[point{x: x, y: y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
