package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runProgram(opcodes []int) error {
	i := 0
	for {
		switch opcodes[i] {
		case 1:
			left := opcodes[opcodes[i+1]]
			right := opcodes[opcodes[i+2]]
			opcodes[opcodes[i+3]] = left + right
			i += 4
		case 2:
			left := opcodes[opcodes[i+1]]
			right := opcodes[opcodes[i+2]]
			opcodes[opcodes[i+3]] = left * right
			i += 4
		case 99:
			return nil
		default:
			return fmt.Errorf("unsupported opcode %d (at position %d)", opcodes[i], i)
		}
	}
}
func TestProgram(t *testing.T) {
	cc := []struct {
		input  []int
		output []int
	}{
		{[]int{1, 0, 0, 0, 99}, []int{2, 0, 0, 0, 99}},
		{[]int{2, 3, 0, 3, 99}, []int{2, 3, 0, 6, 99}},
		{[]int{2, 4, 4, 5, 99, 0}, []int{2, 4, 4, 5, 99, 9801}},
		{[]int{1, 1, 1, 4, 99, 5, 6, 0, 99}, []int{30, 1, 1, 4, 2, 5, 6, 0, 99}},
	}
	for _, c := range cc {
		t.Run("Pogram ", func(t *testing.T) {
			a := assert.New(t)
			err := runProgram(c.input)
			a.NoError(err)
			a.Equal(c.output, c.input)
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
	scanner.Split(SplitByte(','))
	var input []int
	for scanner.Scan() {
		text := scanner.Text()
		opcode, err := strconv.Atoi(strings.TrimSpace(text))
		a.NoError(err)
		input = append(input, opcode)
	}
	a.NoError(scanner.Err())

	input[1] = 12
	input[2] = 2

	err = runProgram(input)
	a.NoError(err)
	a.Equal(3895705, input[0])
}
