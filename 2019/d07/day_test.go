package d07

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

func read1(opcodes []int, i int) int {
	mode := opcodes[i] / 100
	one := opcodes[i+1]
	if mode%10 == 0 { // position mode
		one = opcodes[one]
	}
	return one
}

func read2(opcodes []int, i int) (int, int) {
	mode := opcodes[i] / 100
	one := opcodes[i+1]
	if mode%10 == 0 { // position mode
		one = opcodes[one]
	}

	mode /= 10
	two := opcodes[i+2]
	if mode%10 == 0 { // position mode
		two = opcodes[two]
	}
	return one, two
}

func runProgram(opcodes []int, input <-chan int, output chan<- int) error {
	i := 0
	for {
		switch opcodes[i] % 100 {
		case 1: // add
			left, right := read2(opcodes, i)
			opcodes[opcodes[i+3]] = left + right
			i += 4
		case 2: // multiply
			left, right := read2(opcodes, i)
			opcodes[opcodes[i+3]] = left * right
			i += 4
		case 3: // write
			pos := opcodes[i+1]
			opcodes[pos] = <-input
			i += 2
		case 4: // read
			output <- read1(opcodes, i)
			i += 2
		case 5: // jump-if-true
			left, right := read2(opcodes, i)
			if left != 0 {
				i = right
			} else {
				i += 3
			}
		case 6: // jump-if-false
			left, right := read2(opcodes, i)
			if left == 0 {
				i = right
			} else {
				i += 3
			}
		case 7: // less than
			left, right := read2(opcodes, i)
			if left < right {
				opcodes[opcodes[i+3]] = 1
			} else {
				opcodes[opcodes[i+3]] = 0
			}
			i += 4
		case 8: // equals
			left, right := read2(opcodes, i)
			if left == right {
				opcodes[opcodes[i+3]] = 1
			} else {
				opcodes[opcodes[i+3]] = 0
			}
			i += 4
		case 99: // halt
			close(output)
			return nil
		default:
			return fmt.Errorf("unsupported opcode %d (at position %d)", opcodes[i], i)
		}
	}
}

func runProgramCached(opcodes []int) func(int, int) int {
	type signal struct {
		pahse, input int
	}
	cache := make(map[signal]int)
	return func(phase, input int) int {
		signal := signal{phase, input}
		if v, ok := cache[signal]; ok {
			return v
		}
		c := append([]int(nil), opcodes...)
		in := make(chan int, 2)
		in <- phase
		in <- input
		out := make(chan int, 1)
		err := runProgram(c, in, out)
		if err != nil {
			panic(err)
		}
		v := <-out
		cache[signal] = v
		return v
	}
}

func findMax(opcodes []int) int {
	runner := runProgramCached(opcodes)

	return findMaxRec([]int{0, 1, 2, 3, 4}, 0, runner)
}

func findMaxRec(phases []int, ampli int, runner func(int, int) int) int {
	if len(phases) == 1 {
		return runner(phases[0], ampli)
	}
	max := 0
	for i, p := range phases {
		nextAmpli := runner(p, ampli)
		nextPhases := append([]int(nil), phases[0:i]...)
		nextPhases = append(nextPhases, phases[i+1:]...)
		m := findMaxRec(nextPhases, nextAmpli, runner)
		fmt.Println(nextPhases, m)
		if m > max {
			max = m
		}
	}
	return max
}

func TestProgramMax(t *testing.T) {
	cc := []struct {
		input []int
		max   int
	}{
		{[]int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0}, 43210},
		{[]int{3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23,
			101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0}, 54321},
		{[]int{3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31, -2, 31, 1007, 31, 0, 33,
			1002, 33, 7, 33, 1, 33, 31, 31, 1, 32, 31, 31, 4, 31, 99, 0, 0, 0}, 65210},
	}
	for _, c := range cc {
		t.Run("Pogram ", func(t *testing.T) {
			a := assert.New(t)
			out := findMax(c.input)
			a.Equal(c.max, out)
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

	out := findMax(input)
	a.Equal(880726, out)
}

func possiblePermutations(phases []int) [][]int {
	out := make([][]int, 0)
	if len(phases) == 1 {
		return [][]int{phases}
	}

	for i, v := range phases {
		withoutCurrent := append([]int(nil), phases[0:i]...)
		withoutCurrent = append(withoutCurrent, phases[i+1:]...)
		next := possiblePermutations(withoutCurrent)
		for _, n := range next {
			one := []int{v}
			one = append(one, n...)
			out = append(out, one)
		}
	}
	return out
}

func TestPermutations(t *testing.T) {
	cc := []struct {
		input []int
		n     int
	}{
		{[]int{0}, 1},
		{[]int{1, 2}, 2},
		{[]int{1, 2, 3}, 6},
	}
	for _, c := range cc {
		t.Run("Pogram ", func(t *testing.T) {
			a := assert.New(t)
			p := possiblePermutations(c.input)
			a.Equal(c.n, len(p))
		})
	}
}

func runForever(phases []int, opcodes []int) int {
	firstWire := make(chan int, 2)

	inWire := firstWire

	for _, p := range phases {
		o := append([]int(nil), opcodes...)

		inWire <- p
		outWire := make(chan int, 1)
		go runProgram(o, inWire, outWire)

		inWire = outWire
	}
	firstWire <- 0
	var last int
	for v := range inWire {
		last = v
		firstWire <- v
	}
	return last
}

func findMaxForever(phases []int, opcodes []int) int {
	max := 0
	for _, p := range possiblePermutations(phases) {
		m := runForever(p, opcodes)
		if m > max {
			max = m
		}
	}
	return max
}

func TestProgramForever(t *testing.T) {
	cc := []struct {
		input  []int
		phases []int
		max    int
	}{
		{[]int{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26,
			27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5}, []int{9, 8, 7, 6, 5}, 139629729},
		{[]int{3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54,
			-5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4,
			53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10}, []int{9, 7, 8, 5, 6}, 18216},
	}
	for _, c := range cc {
		t.Run("Pogram ", func(t *testing.T) {
			a := assert.New(t)
			out := runForever(c.phases, c.input)
			a.Equal(c.max, out)
			out = findMaxForever(c.phases, c.input)
			a.Equal(c.max, out)
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

	out := findMaxForever([]int{5, 6, 7, 8, 9}, input)
	a.Equal(4931744, out)
}
