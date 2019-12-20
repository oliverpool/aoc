package d20

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type coord struct {
	x, y int
}

var (
	free = coord{0, 0}
)

type coordMap map[coord]coord

func parseMap(r io.Reader) (coordMap, map[string]coord) {
	paths := make(coordMap)
	partialPortals := make(map[coord]rune)
	portals := make(map[string]coord)

	scanner := bufio.NewScanner(r)
	y := 0
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.ReplaceAll(text, "	", "    ")
		for x, l := range text {
			if l == ' ' || l == '#' {
				continue
			}
			if l == '.' {
				if paths[coord{x, y}] == free {
					paths[coord{x, y}] = free
				}
				continue
			}
			var portal string
			var src coord
			if f, ok := partialPortals[coord{x - 1, y}]; ok {
				portal = string(f) + string(l)

				src = coord{x - 2, y}
				if _, ok := paths[src]; !ok {
					// the paths is not a '.', so take the other side
					src = coord{x + 1, y}
				}
			}
			if f, ok := partialPortals[coord{x, y - 1}]; ok {
				portal = string(f) + string(l)
				src = coord{x, y - 2}
				if _, ok := paths[src]; !ok {
					// the paths is not a '.', so take the other side
					src = coord{x, y + 1}
				}
			}

			if portal != "" {
				if dest, ok := portals[portal]; ok {
					paths[dest] = src
					paths[src] = dest
				} else {
					portals[portal] = src
				}
			}

			partialPortals[coord{x, y}] = l
		}
		y++
	}
	err := scanner.Err()
	if err != nil {
		panic(err)
	}
	return paths, portals
}

func (cm coordMap) String() string {
	var xMin, xMax, yMin, yMax int
	for p := range cm {
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
			s, ok := cm[coord{x, y}]
			if !ok {
				output += " "
			} else if s == free {
				output += "."
			} else {
				output += "@"
			}
		}
		output += "\n"
	}
	return output
}

func (cm coordMap) shortestPath(start, end coord) int {
	current := []coord{start}
	distances := map[coord]int{start: 0}

	for {
		var next []coord
		for _, c := range current {
			d := distances[c]
			neighbors := []coord{
				coord{c.x, c.y - 1},
				coord{c.x + 1, c.y},
				coord{c.x, c.y + 1},
				coord{c.x - 1, c.y},
			}
			dest := cm[c]
			if dest != free {
				neighbors = append(neighbors, dest)
			}
			for _, nei := range neighbors {
				if nei == end {
					return d + 1
				}
				if _, ok := cm[nei]; !ok {
					continue
				}
				if _, ok := distances[nei]; ok {
					continue
				}
				distances[nei] = d + 1
				next = append(next, nei)
			}
		}
		current = next
	}
}

func TestParseMap(t *testing.T) {
	cc := []struct {
		input         string
		free, portals int
		path          int
	}{
		{
			`         A
         A
  #######.#########
  #######.........#
  #######.#######.#
  #######.#######.#
  #######.#######.#
  #####  B    ###.#
BC...##  C    ###.#
  ##.##       ###.#
  ##...DE  F  ###.#
  #####    G  ###.#
  #########.#####.#
DE..#######...###.#
  #.#########.###.#
FG..#########.....#
  ###########.#####
             Z
             Z       `,
			47,
			5,
			23,
		},
		{
			`                   A
                   A
  #################.#############
  #.#...#...................#.#.#
  #.#.#.###.###.###.#########.#.#
  #.#.#.......#...#.....#.#.#...#
  #.#########.###.#####.#.#.###.#
  #.............#.#.....#.......#
  ###.###########.###.#####.#.#.#
  #.....#        A   C    #.#.#.#
  #######        S   P    #####.#
  #.#...#                 #......VT
  #.#.#.#                 #.#####
  #...#.#               YN....#.#
  #.###.#                 #####.#
DI....#.#                 #.....#
  #####.#                 #.###.#
ZZ......#               QG....#..AS
  ###.###                 #######
JO..#.#.#                 #.....#
  #.#.#.#                 ###.#.#
  #...#..DI             BU....#..LF
  #####.#                 #.#####
YN......#               VT..#....QG
  #.###.#                 #.###.#
  #.#...#                 #.....#
  ###.###    J L     J    #.#.###
  #.....#    O F     P    #.#...#
  #.###.#####.#.#####.#####.###.#
  #...#.#.#...#.....#.....#.#...#
  #.#####.###.###.#.#.#########.#
  #...#.#.....#...#.#.#.#.....#.#
  #.###.#####.###.###.#.#.#######
  #.#.........#...#.............#
  #########.###.###.#############
           B   J   C
		   U   P   P               `,
			313, 12, 58,
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			m, portals := parseMap(strings.NewReader(c.input))

			t.Log(m)
			a.Len(m, c.free)
			a.Len(portals, c.portals)
			a.Equal(c.path, m.shortestPath(portals["AA"], portals["ZZ"]))
		})
	}
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	m, portals := parseMap(f)

	a.Equal(692, m.shortestPath(portals["AA"], portals["ZZ"]))
}
