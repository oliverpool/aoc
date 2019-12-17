package d17

import (
	"fmt"
	"io"
	"testing"

	assert "github.com/stretchr/testify/require"
)

type coord struct {
	x, y int
}

// func (c coord) move(d direction) coord {
// 	switch d {
// 	case north:
// 		c.y--
// 	case south:
// 		c.y++
// 	case west:
// 		c.x--
// 	case east:
// 		c.x++
// 	default:
// 		panic(d)
// 	}
// 	return c
// }

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

func TestFirst(t *testing.T) {
	a := assert.New(t)
	var intcodes map[int]int
	var err error
	err = open("./input", func(r io.Reader) error {
		intcodes, err = parseInput(r)
		return err
	})
	a.NoError(err)

	view := make(cameraView)

	pOutput := make(chan int)
	pInput := make(chan int)

	align := 0
	go func() {
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

			if scaLen >= 3 {
				if view[coord{x - 1, y - 1}] == '#' {
					fmt.Println("found", x-1, y)
					align += (x - 1) * y
					view[coord{x - 1, y}] = 'O'
				}
			}
			x++
		}
	}()

	err = runProgram(intcodes, pInput, pOutput)
	a.NoError(err)

	t.Log("\n" + view.String())

	a.Equal(10064, align)
}
