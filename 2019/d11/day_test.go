package d11

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

type coord struct {
	x, y int
}

func (c coord) turnAndGo(d direction) coord {
	switch d {
	case up:
		return coord{c.x, c.y - 1}
	case right:
		return coord{c.x + 1, c.y}
	case down:
		return coord{c.x, c.y + 1}
	case left:
		return coord{c.x - 1, c.y}
	}
	panic("unsupported direction" + strconv.Itoa(int(d)))
}

type color int

const black = color(0)
const white = color(1)

type turn int

const toLeft = turn(0)
const toRight = turn(1)

type direction int

const up = direction(0)
const right = direction(1)
const down = direction(2)
const left = direction(3)

func (d direction) update(t turn) direction {
	switch t {
	case toLeft:
		return direction((d + 3) % 4)
	case toRight:
		return direction((d + 1) % 4)
	}
	panic("unsupported turn" + strconv.Itoa(int(t)))
}

type robot struct {
	Map       map[coord]color
	Coord     coord
	Direction direction
}

func newRobot() robot {
	return robot{
		Map: make(map[coord]color),
	}
}

func (r *robot) run(in <-chan int, out chan<- color) error {
	defer close(out)
	for {
		out <- r.Map[r.Coord]
		ci, ok := <-in
		if !ok {
			return nil
		}
		r.Map[r.Coord] = color(ci)

		r.Direction = r.Direction.update(turn(<-in))

		r.Coord = r.Coord.turnAndGo(r.Direction)
	}
}

func TestProgramSpecial(t *testing.T) {
	cc := []struct {
		input  []int
		output []color
		count  int
	}{
		{[]int{1, 0, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 1, 0}, []color{0, 0, 0, 0, 1, 0, 0, 0}, 6},
	}
	for _, c := range cc {
		t.Run("Program", func(t *testing.T) {
			a := assert.New(t)

			in := make(chan int, len(c.input))
			for _, i := range c.input {
				in <- i
			}
			close(in)

			out := make(chan color, 1)

			robot := newRobot()
			go robot.run(in, out)
			i := 0
			for o := range out {
				a.Equal(c.output[i], o, "output %d", i)
				i++
			}
			a.Equal(len(c.output), i, "output len")

			a.Equal(c.count, len(robot.Map), "map len")
			a.Equal(coord{0, -1}, robot.Coord, "coord")
			a.Equal(left, robot.Direction, "direction")
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

	in := make(chan int)
	outColor := make(chan color)

	out := make(chan int)
	go func() {
		a.NoError(runProgram(input, in, out))
	}()

	go func() {
		for o := range outColor {
			in <- int(o)
		}
		close(in)
	}()

	r := newRobot()
	r.run(out, outColor)

	a.Equal(1709, len(r.Map))
}
