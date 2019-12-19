package d19

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		for y := 0; y < h; y++ {
			go func(x, y int) {
				o, err := getBeam(intcodes, x, y)
				a.NoError(err)
				output <- o
			}(x, y)
		}
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
