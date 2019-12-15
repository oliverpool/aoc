package d09

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

func read1(opcodes map[int]int, i, base int) int {
	mode := opcodes[i] / 100
	one := opcodes[i+1]
	if mode%10 == 0 { // position mode
		one = opcodes[one]
	} else if mode%10 == 2 { // base mode
		one = opcodes[one+base]
	}
	return one
}

func read2(opcodes map[int]int, i, base int) (int, int) {
	mode := opcodes[i] / 100
	one := opcodes[i+1]
	if mode%10 == 0 { // position mode
		one = opcodes[one]
	} else if mode%10 == 2 { // base mode
		one = opcodes[one+base]
	}

	mode /= 10
	two := opcodes[i+2]
	if mode%10 == 0 { // position mode
		two = opcodes[two]
	} else if mode%10 == 2 { // base mode
		two = opcodes[two+base]
	}
	return one, two
}

func readDest(opcodes map[int]int, i, delta, base int) int {
	dest := opcodes[i+delta]
	mode := opcodes[i] / 10
	for delta > 0 {
		mode /= 10
		delta--
	}
	if mode%10 == 2 { // base mode
		dest += base
	}
	return dest
}

func runProgram(opcodes map[int]int, input <-chan int, output chan<- int) error {
	i := 0
	base := 0
	for {
		switch opcodes[i] % 100 {
		case 1: // add
			left, right := read2(opcodes, i, base)
			dest := readDest(opcodes, i, 3, base)
			opcodes[dest] = left + right
			i += 4
		case 2: // multiply
			left, right := read2(opcodes, i, base)
			dest := readDest(opcodes, i, 3, base)
			opcodes[dest] = left * right
			i += 4
		case 3: // write
			dest := readDest(opcodes, i, 1, base)
			opcodes[dest] = <-input
			i += 2
		case 4: // read
			output <- read1(opcodes, i, base)
			i += 2
		case 5: // jump-if-true
			left, right := read2(opcodes, i, base)
			if left != 0 {
				i = right
			} else {
				i += 3
			}
		case 6: // jump-if-false
			left, right := read2(opcodes, i, base)
			if left == 0 {
				i = right
			} else {
				i += 3
			}
		case 7: // less than
			left, right := read2(opcodes, i, base)
			dest := readDest(opcodes, i, 3, base)
			if left < right {
				opcodes[dest] = 1
			} else {
				opcodes[dest] = 0
			}
			i += 4
		case 8: // equals
			left, right := read2(opcodes, i, base)
			dest := readDest(opcodes, i, 3, base)
			if left == right {
				opcodes[dest] = 1
			} else {
				opcodes[dest] = 0
			}
			i += 4
		case 9: // adjust relative base
			adj := read1(opcodes, i, base)
			base += adj
			i += 2
		case 99: // halt
			close(output)
			return nil
		default:
			close(output)
			return fmt.Errorf("unsupported opcode %d (at position %d)", opcodes[i], i)
		}
	}
}

func TestProgramSpecial(t *testing.T) {
	cc := []struct {
		input []int
		out   []int
	}{
		{[]int{104, 1125899906842624, 99}, []int{1125899906842624}},
		{[]int{1102, 34915192, 34915192, 7, 4, 7, 99, 0}, []int{1219070632396864}},
		{[]int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
			[]int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99}},
	}
	for _, c := range cc {
		t.Run("Program ", func(t *testing.T) {
			a := assert.New(t)
			input := make(map[int]int, len(c.input))
			for i, v := range c.input {
				input[i] = v
			}
			outc := make(chan int, len(c.input))
			a.NoError(runProgram(input, nil, outc))
			var out []int
			for o := range outc {
				out = append(out, o)
			}
			a.Equal(c.out, out)
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
	input := make(map[int]int)
	i := 0
	for scanner.Scan() {
		text := scanner.Text()
		opcode, err := strconv.Atoi(strings.TrimSpace(text))
		a.NoError(err)
		input[i] = opcode
		i++
	}
	a.NoError(scanner.Err())

	in := make(chan int, 1)
	in <- 1
	out := make(chan int)
	go func() {
		a.NoError(runProgram(input, in, out))
	}()

	var o int
	for o = range out {
		t.Log(o)
	}
	a.Equal(3638931938, o)
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(SplitByte(','))
	input := make(map[int]int)
	i := 0
	for scanner.Scan() {
		text := scanner.Text()
		opcode, err := strconv.Atoi(strings.TrimSpace(text))
		a.NoError(err)
		input[i] = opcode
		i++
	}
	a.NoError(scanner.Err())

	in := make(chan int, 1)
	in <- 2
	out := make(chan int)
	go func() {
		a.NoError(runProgram(input, in, out))
	}()

	var o int
	for o = range out {
		t.Log(o)
	}
	a.Equal(86025, o)
}
