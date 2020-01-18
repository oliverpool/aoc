package d19

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type coord struct{ x, y int }

var beamCache = make(map[coord]bool)

func mustGetBeam(intcodes map[int]int, x, y int) bool {
	if v, ok := beamCache[coord{x, y}]; ok {
		return v
	}
	o, err := getBeam(intcodes, x, y)
	if err != nil {
		panic(err)
	}
	beamCache[coord{x, y}] = o > 0
	return o > 0
}

func getBeam(intcodes map[int]int, x, y int) (int, error) {
	copy := make(map[int]int)
	for i, v := range intcodes {
		copy[i] = v
	}
	pInput := make(chan int, 2)
	pInput <- x
	pInput <- y
	pOutput := make(chan int, 1)
	err := runProgram(copy, pInput, pOutput)
	return <-pOutput, err
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	var intcodes map[int]int
	var err error
	err = open("./input", func(r io.Reader) error {
		intcodes, err = parseInput(r)
		return err
	})
	a.NoError(err)

	w, h := 50, 50

	output := make(chan int, 10)
	for x := 0; x < w; x++ {
		go func(x int) {
			for y := 0; y < h; y++ {
				o, err := getBeam(intcodes, x, y)
				a.NoError(err)
				output <- o
			}
		}(x)
	}
	n := 0
	area := 0
	for o := range output {
		area++
		if o > 0 {
			n++
		}
		if area >= w*h {
			break
		}
	}

	a.Equal(223, n)
}

/*
w,h=100,100
find x,y, such that
(x+w,y) == 1
(x+w,y-1) == 0

(x,y+h) == 1
(x,y+h-1) == 0

x->
y
|
V

a*x + b*y = 0 for both
*/

type params struct {
	a, b int
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	var intcodes map[int]int
	var err error
	err = open("./input", func(r io.Reader) error {
		intcodes, err = parseInput(r)
		return err
	})
	a.NoError(err)

	initial := 100

	// Find a & b of: a*x + b*y = 0
	x := initial
	var pr params
	for y := 0; true; y++ {
		o, err := getBeam(intcodes, x, y)
		a.NoError(err)
		if o > 0 {
			pr = params{y, -initial}
			break
		}
	}

	y := initial
	var pl params
	for x := 0; true; x++ {
		o, err := getBeam(intcodes, x, y)
		a.NoError(err)
		if o > 0 {
			pl = params{-initial, x}
			break
		}
	}
	t.Log(pr, pl)
	/*
		Now find approx x and y, so that:
			al*x + bl*(y+99) = 0
			ar*(x+99) + br*y = 0

			br*al*x + bl*(br*y+br*99) = 0

			br*al*x + bl*(br*99-ar*(x+99)) = 0
			(br*al-bl*ar)*x + bl*br*99-bl*ar*99 = 0

			x = bl*(ar*99 - br*99)/(br*al-bl*ar)
			=> y = - ar*(x+99)/br
	*/

	w, h := 99, 99
	x = pl.b * (pr.a*w - pr.b*h) / (pr.b*pl.a - pl.b*pr.a)
	y = -pr.a * (x + w) / pr.b

	topRight := func() bool {
		return mustGetBeam(intcodes, x+w, y)
	}
	bottomLeft := func() bool {
		return mustGetBeam(intcodes, x, y+h)
	}

	a.True(topRight())
	a.True(bottomLeft())

	xPrev, yPrev := 0, 0
	for xPrev != x || yPrev != y {
		xPrev, yPrev = x, y

		// go to the lowest y
		for mustGetBeam(intcodes, x+w, y-1) {
			y--
		}

		// go to the lowest x
		for mustGetBeam(intcodes, x-1, y+h) {
			x--
		}

		if xPrev != x || yPrev != y {
			continue
		}

		// attempt a jump
		x--
		y--

		// go to the lowest y
		for mustGetBeam(intcodes, x+w, y-1) {
			y--
		}

		// go to the lowest x
		for mustGetBeam(intcodes, x-1, y+h) {
			x--
		}

		// go to the lowest y
		for mustGetBeam(intcodes, x+w, y-1) {
			y--
		}

		if !topRight() || !bottomLeft() {
			x, y = xPrev, yPrev
		}
	}

	a.Equal(9480761, 10000*x+y)
}
