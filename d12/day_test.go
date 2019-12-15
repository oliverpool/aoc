package d12

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func sign(a, b int) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func abs(v int) int {
	if v > 0 {
		return v
	} else {
		return -v
	}
}

func simulateMotion(pos, vel []int, n int) ([]int, []int) {
	for n > 0 {
		// compute new velocity
		for i := range pos {
			for j := i + 1; j < len(pos); j++ {
				s := sign(pos[i], pos[j])
				vel[i] -= s
				vel[j] += s
			}
		}

		// update position
		for i, v := range vel {
			pos[i] += v
		}
		n--
	}

	return pos, vel
}

type point struct {
	x, y, z int
}

func kineticAfter(moons []point, n int) int {
	var pos [3][]int
	var vel [3][]int
	for _, m := range moons {
		pos[0] = append(pos[0], m.x)
		pos[1] = append(pos[1], m.y)
		pos[2] = append(pos[2], m.z)

		vel[0] = append(vel[0], 0)
		vel[1] = append(vel[1], 0)
		vel[2] = append(vel[2], 0)
	}
	for i := range pos {
		pos[i], vel[i] = simulateMotion(pos[i], vel[i], n)
	}

	var sum int
	for i := range pos[0] {
		var ps int
		ps += abs(pos[0][i])
		ps += abs(pos[1][i])
		ps += abs(pos[2][i])

		var vs int
		vs += abs(vel[0][i])
		vs += abs(vel[1][i])
		vs += abs(vel[2][i])

		sum += ps * vs
	}

	return sum
}

func TestEvolution(t *testing.T) {
	cc := []struct {
		position    []int
		steps       int
		positionEnd []int
		velocityEnd []int
	}{
		{
			[]int{-1, 2, 4, 3},
			10,
			[]int{2, 1, 3, 2},
			[]int{-3, -1, 3, 1},
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			vel := make([]int, len(c.position))
			pos, vel := simulateMotion(c.position, vel, c.steps)
			a.Equal(c.positionEnd, pos)
			a.Equal(c.velocityEnd, vel)
		})
	}
}

func TestKinetic(t *testing.T) {
	cc := []struct {
		moons   []point
		steps   int
		kinetic int
	}{
		{
			[]point{{-1, 0, 2}, {2, -10, -7}, {4, -8, 8}, {3, 5, -1}},
			10,
			179,
		},
		{
			[]point{{-8, -10, 0}, {5, 5, 10}, {2, -7, 3}, {9, -8, -3}},
			100,
			1940,
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			k := kineticAfter(c.moons, c.steps)
			a.Equal(c.kinetic, k)
		})
	}
}

func TestParse(t *testing.T) {
	cc := []struct {
		input string
		moons []point
	}{
		{
			`<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>`,
			[]point{{-1, 0, 2}, {2, -10, -7}, {4, -8, 8}, {3, 5, -1}},
		},
		{
			`<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>
`,
			[]point{{-8, -10, 0}, {5, 5, 10}, {2, -7, 3}, {9, -8, -3}},
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			m, err := parseMoons(strings.NewReader(c.input))
			a.NoError(err)
			a.Equal(c.moons, m)

		})
	}
}

func parseMoons(r io.Reader) ([]point, error) {
	var moons []point
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		var p point
		_, err := fmt.Sscanf(text, "<x=%d, y=%d, z=%d>", &p.x, &p.y, &p.z)
		moons = append(moons, p)
		if err != nil {
			return moons, err
		}
	}
	return moons, scanner.Err()
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	moons, err := parseMoons(f)
	f.Close()
	a.NoError(err)

	k := kineticAfter(moons, 1000)
	a.Equal(12082, k)
}

func repeatingMotion(pos []int) int {
	initialPos := make([]int, len(pos))
	copy(initialPos, pos)
	vel := make([]int, len(pos))
	n := 0

	isBackToInial := func() bool {
		for _, v := range vel {
			if v != 0 {
				return false
			}
		}
		for i, p := range pos {
			if p != initialPos[i] {
				return false
			}
		}
		return true
	}

	for {
		// compute new velocity
		for i := range pos {
			for j := i + 1; j < len(pos); j++ {
				s := sign(pos[i], pos[j])
				vel[i] -= s
				vel[j] += s
			}
		}

		// update position
		for i, v := range vel {
			pos[i] += v
		}
		n++
		if isBackToInial() {
			return n
		}
	}
}

func backToInial(moons []point) int {
	var pos [3][]int
	var steps []int
	for _, m := range moons {
		pos[0] = append(pos[0], m.x)
		pos[1] = append(pos[1], m.y)
		pos[2] = append(pos[2], m.z)
	}
	for i := range pos {
		steps = append(steps, repeatingMotion(pos[i]))
	}
	return lcm(steps...)
}

func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	if x < 0 {
		return -x
	}
	return x
}

func lcm(i ...int) int {
	if len(i) < 1 {
		panic("at least one argument is necessary")
	} else if len(i) == 1 {
		return i[0]
	} else {
		a, b := i[0], lcm(i[1:]...)
		return a * b / gcd(a, b)
	}
}

func TestRepeatingMotion(t *testing.T) {
	cc := []struct {
		coord []int
		steps int
	}{
		{
			[]int{-1, 2, 4, 3},
			18,
		},
		{
			[]int{0, -10, -8, 5},
			28,
		},
		{
			[]int{2, -7, 8, -1},
			44,
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			n := repeatingMotion(c.coord)
			a.Equal(c.steps, n)
		})
	}

	a := assert.New(t)
	a.Equal(2772, lcm(18, 28, 44))
}

func TestBackToInitial(t *testing.T) {
	cc := []struct {
		moons []point
		steps int
	}{
		{
			[]point{{-1, 0, 2}, {2, -10, -7}, {4, -8, 8}, {3, 5, -1}},
			2772,
		},
		{
			[]point{{-8, -10, 0}, {5, 5, 10}, {2, -7, 3}, {9, -8, -3}},
			4686774924,
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			n := backToInial(c.moons)
			a.Equal(c.steps, n)
		})
	}
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	moons, err := parseMoons(f)
	f.Close()
	a.NoError(err)

	n := backToInial(moons)
	a.Equal(295693702908636, n)
}
