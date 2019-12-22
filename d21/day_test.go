package d21

import (
	"fmt"
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

	output := make(chan int, 10)
	input := make(chan int, 10)

	go func() {
		err := runProgram(intcodes, input, output)
		a.NoError(err)
	}()

	go func() {
		// jump if A,B or C is a hole AND D is a ground
		for _, l := range `NOT A J
NOT B T
OR T J
NOT C T
OR T J
AND D J
WALK
` {
			input <- int(l)
		}
	}()

	var last int
	for o := range output {
		last = o
		fmt.Print(string(o))
	}

	a.Equal(19350938, last)
}
