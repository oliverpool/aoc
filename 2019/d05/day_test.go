package d05

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

func runProgram(opcodes []int, input int) ([]int, error) {
	i := 0
	var output []int
	for {
		switch opcodes[i] % 100 {
		case 1:
			mode := opcodes[i] / 100
			left := opcodes[i+1]
			if mode%10 == 0 { // position mode
				left = opcodes[left]
			}
			mode /= 10
			right := opcodes[i+2]
			if mode%10 == 0 { // position mode
				right = opcodes[right]
			}
			opcodes[opcodes[i+3]] = left + right
			i += 4
		case 2:
			mode := opcodes[i] / 100
			left := opcodes[i+1]
			if mode%10 == 0 { // position mode
				left = opcodes[left]
			}
			mode /= 10
			right := opcodes[i+2]
			if mode%10 == 0 { // position mode
				right = opcodes[right]
			}
			opcodes[opcodes[i+3]] = left * right
			i += 4
		case 3:
			pos := opcodes[i+1]
			opcodes[pos] = input
			i += 2
		case 4:
			mode := opcodes[i] / 100
			pos := opcodes[i+1]
			if mode%10 == 0 { // position mode
				pos = opcodes[pos]
			}
			output = append(output, pos)
			i += 2
		case 99:
			return output, nil
		default:
			return nil, fmt.Errorf("unsupported opcode %d (at position %d)", opcodes[i], i)
		}
	}
}

func TestProgram(t *testing.T) {
	cc := []struct {
		input  []int
		output []int
		sdtin  int
		sdtout []int
	}{
		{[]int{1, 0, 0, 0, 99}, []int{2, 0, 0, 0, 99}, 0, nil},
		{[]int{2, 3, 0, 3, 99}, []int{2, 3, 0, 6, 99}, 0, nil},
		{[]int{2, 4, 4, 5, 99, 0}, []int{2, 4, 4, 5, 99, 9801}, 0, nil},
		{[]int{1, 1, 1, 4, 99, 5, 6, 0, 99}, []int{30, 1, 1, 4, 2, 5, 6, 0, 99}, 0, nil},
		{[]int{3, 0, 4, 0, 99}, []int{1, 0, 4, 0, 99}, 1, []int{1}},
		{[]int{3, 0, 104, 42, 99}, []int{1, 0, 104, 42, 99}, 1, []int{42}},
	}
	for _, c := range cc {
		t.Run("Pogram ", func(t *testing.T) {
			a := assert.New(t)
			out, err := runProgram(c.input, c.sdtin)
			a.NoError(err)
			a.Equal(c.output, c.input)
			a.Equal(c.sdtout, out)
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

	out, err := runProgram(input, 1)
	a.NoError(err)
	a.Equal(3895705, out[len(out)-1])
}
