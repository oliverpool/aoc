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
		case 5:
			mode := opcodes[i] / 100
			first := opcodes[i+1]
			if mode%10 == 0 { // position mode
				first = opcodes[first]
			}
			mode /= 10
			second := opcodes[i+2]
			if mode%10 == 0 { // position mode
				second = opcodes[second]
			}
			if first != 0 {
				i = second
			} else {
				i += 3
			}
		case 6:
			mode := opcodes[i] / 100
			first := opcodes[i+1]
			if mode%10 == 0 { // position mode
				first = opcodes[first]
			}
			mode /= 10
			second := opcodes[i+2]
			if mode%10 == 0 { // position mode
				second = opcodes[second]
			}
			if first == 0 {
				i = second
			} else {
				i += 3
			}
		case 7:
			mode := opcodes[i] / 100
			first := opcodes[i+1]
			if mode%10 == 0 { // position mode
				first = opcodes[first]
			}
			mode /= 10
			second := opcodes[i+2]
			if mode%10 == 0 { // position mode
				second = opcodes[second]
			}
			if first < second {
				opcodes[opcodes[i+3]] = 1
			} else {
				opcodes[opcodes[i+3]] = 0
			}
			i += 4
		case 8:
			mode := opcodes[i] / 100
			first := opcodes[i+1]
			if mode%10 == 0 { // position mode
				first = opcodes[first]
			}
			mode /= 10
			second := opcodes[i+2]
			if mode%10 == 0 { // position mode
				second = opcodes[second]
			}
			if first == second {
				opcodes[opcodes[i+3]] = 1
			} else {
				opcodes[opcodes[i+3]] = 0
			}
			i += 4
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
	a.Equal(4511442, out[len(out)-1])
}

func TestProgramJump(t *testing.T) {
	longCase := []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
		1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
		999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99}
	cc := []struct {
		input  []int
		sdtin  int
		sdtout int
	}{
		{[]int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}, 7, 0},
		{[]int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}, 8, 1},
		{[]int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}, 9, 0},

		{[]int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8}, 7, 1},
		{[]int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8}, 8, 0},
		{[]int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8}, 9, 0},

		{[]int{3, 3, 1108, -1, 8, 3, 4, 3, 99}, 7, 0},
		{[]int{3, 3, 1108, -1, 8, 3, 4, 3, 99}, 8, 1},
		{[]int{3, 3, 1108, -1, 8, 3, 4, 3, 99}, 9, 0},

		{[]int{3, 3, 1107, -1, 8, 3, 4, 3, 99}, 7, 1},
		{[]int{3, 3, 1107, -1, 8, 3, 4, 3, 99}, 8, 0},
		{[]int{3, 3, 1107, -1, 8, 3, 4, 3, 99}, 9, 0},

		{[]int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9}, 0, 0},
		{[]int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9}, 1, 1},
		{[]int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9}, 2, 1},

		{[]int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1}, 0, 0},
		{[]int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1}, 1, 1},
		{[]int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1}, 2, 1},

		{longCase, 6, 999},
		{longCase, 7, 999},
		{longCase, 8, 1000},
		{longCase, 9, 1001},
	}
	for _, c := range cc {
		t.Run("Pogram ", func(t *testing.T) {
			a := assert.New(t)
			out, err := runProgram(c.input, c.sdtin)
			a.NoError(err)
			a.Equal(c.sdtout, out[len(out)-1])
		})
	}
}

func TestSecond(t *testing.T) {
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

	out, err := runProgram(input, 5)
	a.NoError(err)
	a.Equal(12648139, out[len(out)-1])
}
