package d13

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func open(path string, cb func(io.Reader) error) error {
	f, err := os.Open("./input")
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

type point struct {
	x, y int
}
type tile int

func (t tile) String() string {
	switch t {
	case empty:
		return " "
	case wall:
		return "#"
	case block:
		return "*"
	case paddle:
		return "="
	case ball:
		return "o"
	}
	return strconv.Itoa(int(t))
}

const (
	empty  = tile(0)
	wall   = tile(1)
	block  = tile(2)
	paddle = tile(3)
	ball   = tile(4)
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

func parseTiles(intcodes <-chan int) map[point]tile {
	tiles := make(map[point]tile)
	for {
		x, ok := <-intcodes
		if !ok {
			return tiles
		}
		y := <-intcodes
		t := tile(<-intcodes)
		tiles[point{x, y}] = t
	}
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

	pOutput := make(chan int)
	go func() {
		err := runProgram(intcodes, nil, pOutput)
		a.NoError(err)
	}()

	tiles := parseTiles(pOutput)
	count := 0
	for _, t := range tiles {
		if t == block {
			count++
		}
	}
	a.Equal(205, count)
}

func gameToString(tiles map[point]tile) string {
	var xMin, xMax, yMin, yMax int
	for p := range tiles {
		if p.x < xMin {
			xMin = p.x
		} else if p.x > xMax {
			xMax = p.x
		}
		if p.y < yMin {
			yMin = p.y
		} else if p.y > yMax {
			yMax = p.y
		}
	}

	output := ""
	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			output += fmt.Sprint(tiles[point{x, y}])
		}
		output += "\n"
	}
	return output
}

func play(intcodes <-chan int, joystick chan<- int) int {
	tiles := make(map[point]tile)

	var pBall, pPaddle point
	pBallDirection := 0
	nextMove := 0
	for {
		select {
		case x, ok := <-intcodes:
			if !ok {
				return int(tiles[point{-1, 0}])
			}
			y := <-intcodes
			t := tile(<-intcodes)
			p := point{x, y}
			tiles[p] = t
			if t == ball {
				if pBall.x > p.x {
					pBallDirection = -1
				} else if pBall.x == p.x {
					pBallDirection = 0
				} else {
					pBallDirection = 1
				}
				pBall = p
			} else if t == paddle {
				pPaddle = p
			}
			if pBall.x == pPaddle.x {
				if pBall.y == pPaddle.y-1 {
					nextMove = 0
				} else {
					nextMove = pBallDirection
				}
			} else if pBall.x+pBallDirection > pPaddle.x {
				nextMove = 1
			} else if pBall.x+pBallDirection < pPaddle.x {
				nextMove = -1
			} else {
				nextMove = 0
			}

		case joystick <- nextMove:
			// fmt.Println(gameToString(tiles))
		}
	}
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	var intcodes map[int]int
	var err error
	err = open("./input", func(r io.Reader) error {
		intcodes, err = parseInput(r)
		return err
	})
	a.NoError(err)

	// put quarters
	intcodes[0] = 2

	pOutput := make(chan int)
	pInput := make(chan int)
	go func() {
		err := runProgram(intcodes, pInput, pOutput)
		a.NoError(err)
	}()

	score := play(pOutput, pInput)

	a.Equal(10292, score)
}
