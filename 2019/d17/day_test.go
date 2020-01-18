package d17

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type coord struct {
	x, y int
}

func (c coord) move(d direction) coord {
	switch d {
	case north:
		c.y--
	case south:
		c.y++
	case west:
		c.x--
	case east:
		c.x++
	default:
		panic(d)
	}
	return c
}

type pixel byte

type cameraView map[coord]pixel

func (sm cameraView) String() string {
	var xMin, xMax, yMin, yMax int
	for p := range sm {
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
			s, ok := sm[coord{x, y}]
			if !ok {
				output += " "
			} else {
				output += string(s)
			}
		}
		output += "\n"
	}
	return output
}

type direction int

const (
	none  = direction(0)
	north = direction(1)
	east  = direction(2)
	south = direction(3)
	west  = direction(4)
)

func getDirection(p pixel) direction {
	switch p {
	case '^':
		return north
	case '>':
		return east
	case 'v':
		return south
	case '<':
		return west
	default:
		panic("unknown direction: " + string(p))
	}
}

func (d direction) others() []direction {
	all := []direction{north, east, south, west}
	for i := 0; i < len(all); i++ {
		if all[i]%2 != d%2 {
			// keep
			continue
		}
		// remove
		all = append(all[:i], all[i+1:]...)
	}
	return all
}

func (d direction) turn(next direction) string {
	if d == next {
		panic("no turn")
	}
	if d%2 == next%2 {
		panic("u turn")
	}
	if (d%4)+1 == next {
		return "R"
	}
	return "L"
}

func (cv cameraView) computePath(currentC coord) []string {
	currentD := getDirection(cv[currentC])
	var nextD direction
	nextC := coord{-1, -1}
	var path []string

	for {
		for _, d := range currentD.others() {
			moved := currentC.move(d)
			if cv[moved] != '#' {
				continue
			}
			nextC = moved
			nextD = d
			break
		}
		if nextD != currentD {
			path = append(path, currentD.turn(nextD))
			currentD = nextD
		}

		previousC := nextC
		l := 0
		for cv[nextC] == '#' {
			previousC = nextC
			nextC = nextC.move(currentD)
			l++
		}
		if l == 0 {
			return path
		}
		path = append(path, strconv.Itoa(l))
		currentC = previousC

	}
}

func getCameraView(intcodes map[int]int) (cameraView, int, coord) {
	view := make(cameraView)

	pOutput := make(chan int)
	pInput := make(chan int)
	go runProgram(intcodes, pInput, pOutput)

	align := 0
	var robot coord

	x, y := 0, 0
	scaLen := 0 // how many pixels are scaffold
	for s := range pOutput {
		b := pixel(s)

		if b == '#' {
			scaLen++
		} else {
			scaLen = 0
		}
		if b == '\n' {
			x = 0
			y++
			continue
		}

		view[coord{x, y}] = b

		if bytes.ContainsRune([]byte{'^', '>', 'v', '<'}, rune(b)) {
			robot = coord{x, y}
		}

		if scaLen >= 3 {
			if view[coord{x - 1, y - 1}] == '#' {
				// fmt.Println("found", x-1, y)
				// view[coord{x - 1, y}] = 'O'
				align += (x - 1) * y
			}
		}
		x++
	}
	return view, align, robot
}

func runRobot(intcodes map[int]int, input string) (cameraView, int) {
	intcodes[0] = 2
	pOutput := make(chan int)
	pInput := make(chan int, len(input)+2)
	go runProgram(intcodes, pInput, pOutput)

	for _, b := range input {
		pInput <- int(b)
	}
	pInput <- int('n')
	pInput <- int('\n')

	view := make(cameraView)
	x, y := 0, 0
	var last int
	for o := range pOutput {
		last = o
		b := pixel(o)

		if b == '\n' {
			x = 0
			y++
			continue
		}

		view[coord{x, y}] = b
		x++
	}
	return view, last
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

	view, align, _ := getCameraView(intcodes)

	t.Log("\n" + view.String())

	a.Equal(10064, align)
}

func TestSecondPre(t *testing.T) {
	a := assert.New(t)
	var intcodes map[int]int
	var err error
	err = open("./input", func(r io.Reader) error {
		intcodes, err = parseInput(r)
		return err
	})
	a.NoError(err)

	intcodes2 := make(map[int]int, len(intcodes))
	for i, v := range intcodes {
		intcodes2[i] = v
	}

	view, _, start := getCameraView(intcodes)

	t.Log("\n" + view.String())

	a.Equal(coord{24, 0}, start)
	a.Equal("^", string(view[start]))

	path := view.computePath(start)
	a.Equal(78, len(path))
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

	input := strings.Join([]string{"A,A,B,C,B,C,B,C,B,A\n",
		"L,10,L,8,R,8,L,8,R,6\n",
		"R,6,R,8,R,8\n",
		"R,6,R,6,L,8,L,10\n",
	}, "")

	view, dust := runRobot(intcodes, input)
	t.Log("\n" + view.String())

	a.Equal(1197725, dust)
}
