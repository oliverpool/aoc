package d25

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
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

func open(path string, cb func(io.Reader) error) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return cb(f)
}

func parseInput(r io.Reader) (map[int]int, error) {
	splitByte := func(sep byte) bufio.SplitFunc {
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

	intcodes := make(map[int]int)
	i := 0

	scanner := bufio.NewScanner(r)
	scanner.Split(splitByte(','))
	for scanner.Scan() {
		text := scanner.Text()

		var v int
		_, err := fmt.Sscanf(text, "%d", &v)
		if err != nil {
			return intcodes, err
		}
		intcodes[i] = v
		i++
	}
	return intcodes, scanner.Err()
}
