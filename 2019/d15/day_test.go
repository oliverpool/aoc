package d15

import (
	"fmt"
	"io"
	"testing"

	assert "github.com/stretchr/testify/require"
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

type statusMap map[coord]status

func (sm statusMap) moveToUnkown(c coord) (direction, coord) {
	directions := []direction{
		north,
		east,
		south,
		west,
	}
	for _, d := range directions {
		moved := c.move(d)
		if _, ok := sm[moved]; !ok {
			return d, moved
		}
	}
	return none, coord{}
}

func (sm statusMap) String() string {
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
			if x == 0 && y == 0 {
				output += "X"
			} else if !ok {
				output += " "
			} else if s == free {
				output += "·"
			} else if s == wall {
				output += "█"
			} else if s == oxygen {
				output += "O"
			} else {
				output += fmt.Sprint(s)
			}
		}
		output += "\n"
	}
	return output
}

type distanceMap map[coord]int

func (dm distanceMap) moveLowest(c coord) coord {
	directions := []direction{
		north,
		south,
		west,
		east,
	}
	remaining := dm[c] - 1
	var next coord
	for _, d := range directions {
		moved := c.move(d)
		if d, ok := dm[moved]; ok && d <= remaining {
			remaining = dm[moved]
			next = moved
		}
	}
	return next
}

func searchOxygen(statuses <-chan status, directions chan<- direction) int {
	var current coord
	m := statusMap{current: free}
	var breadcrumb []direction
	distances := make(distanceMap)
	var oxygenCoord coord
	for {
		d := none
		var moved coord
		for d == none {
			d, moved = m.moveToUnkown(current)
			if d != none {
				break
			}

			if len(breadcrumb) == 0 {
				fmt.Println(m)
				return pathLengthToOrigin(oxygenCoord, distances)
			}

			l := len(breadcrumb) - 1
			d, breadcrumb = breadcrumb[l], breadcrumb[:l]

			directions <- d
			<-statuses
			current = current.move(d)
			d, moved = m.moveToUnkown(current)
		}
		directions <- d
		s := <-statuses
		m[moved] = s
		if s != wall {
			current = moved
			breadcrumb = append(breadcrumb, d.back())
			distances[current] = len(breadcrumb)
			if s == oxygen {
				fmt.Println("oxygen", moved)
				oxygenCoord = current
			}
		}
	}
}

func pathLengthToOrigin(start coord, distances distanceMap) int {
	distance := 0
	var end coord
	for start != end {
		start = distances.moveLowest(start)
		distance++
	}

	return distance
}

type status int

const (
	wall   = status(0)
	free   = status(1)
	oxygen = status(2)
)

type direction int

const (
	none  = direction(0)
	north = direction(1)
	south = direction(2)
	west  = direction(3)
	east  = direction(4)
)

func (d direction) back() direction {
	switch d {
	case north:
		return south
	case south:
		return north
	case east:
		return west
	case west:
		return east
	default:
		panic(d)
	}
}

func trainingMap(m map[coord]status) func(statuses chan<- status, directions <-chan direction) {
	var current coord
	return func(statuses chan<- status, directions <-chan direction) {
		for d := range directions {
			moved := current.move(d)
			s, ok := m[moved]
			if !ok {
				statuses <- wall
			} else {
				if s != wall {
					current = moved
				}
				statuses <- s
			}
		}
	}
}

func TestMap(t *testing.T) {
	cc := []struct {
		m        statusMap
		distance int
	}{
		{
			statusMap{
				{0, 0}:  free,
				{1, 0}:  free,
				{0, 1}:  free,
				{-1, 1}: oxygen,
			},
			2,
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)

			statuses := make(chan status)
			directions := make(chan direction)

			go trainingMap(c.m)(statuses, directions)
			distance := searchOxygen(statuses, directions)
			a.Equal(c.distance, distance)
		})
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

	statuses := make(chan status)
	directions := make(chan direction)

	pOutput := make(chan int)
	pInput := make(chan int)

	go func() {
		for d := range directions {
			pInput <- int(d)
		}
	}()

	go func() {
		for s := range pOutput {
			statuses <- status(s)
		}
	}()

	go func() {
		err := runProgram(intcodes, pInput, pOutput)
		a.NoError(err)
	}()

	distance := searchOxygen(statuses, directions)

	a.Equal(262, distance)
}
