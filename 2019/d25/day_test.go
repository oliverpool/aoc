package d25

import (
	"bufio"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runComputer(intcodes map[int]int, r io.Reader, w io.Writer) error {
	pOutput := make(chan int, 10)
	go func() {
		for b := range pOutput {
			w.Write([]byte{byte(b)})
		}
	}()

	pInput := make(chan int, 10)
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			text := scanner.Text()
			if expanded, ok := map[string]string{
				"n": "north",
				"s": "south",
				"w": "west",
				"e": "east",
			}[text]; ok {
				text = expanded
			}
			for _, b := range text {
				pInput <- int(b)
			}
			pInput <- '\n'
		}
		err := scanner.Err()
		if err != nil {
			panic(err)
		}
	}()

	return runProgram(intcodes, pInput, pOutput)
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

	t.Skip()

	// go test -c && ./d25.test
	/*
		- ornament
		- astrolabe
		- sand
	*/
	runComputer(intcodes, os.Stdin, os.Stdout)
}
