package d22

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type process []Shuffler

func (p process) NewPos(n int64) int64 {
	for _, s := range p {
		n = s.NewPos(n)
	}
	return n
}

func (p process) Inverse() Shuffler {
	inv := make(process, len(p))
	for i, s := range p {
		inv[len(p)-i-1] = s.Inverse()
	}
	return inv
}

type Shuffler interface {
	NewPos(int64) int64
	Inverse() Shuffler
}

type newStack struct {
	len int64
}

func (ns newStack) NewPos(i int64) int64 {
	return ns.len - i - 1
}

func (ns newStack) Inverse() Shuffler {
	return ns
}

type cutN struct {
	n, len int64
}

func (cn cutN) NewPos(i int64) int64 {
	return (((i - cn.n) % cn.len) + cn.len) % cn.len
}

func (cn cutN) Inverse() Shuffler {
	return cutN{cn.len - cn.n, cn.len}
}

type incN struct {
	n, len int64
}

func (in incN) NewPos(i int64) int64 {

	return (i * in.n) % in.len
}

func (in incN) Inverse() Shuffler {
	r, u, v := in.n, int64(1), int64(0)
	r2, u2, v2 := in.len, int64(0), int64(1)
	for r2 != 0 {
		q := r / r2
		r, u, v, r2, u2, v2 = r2, u2, v2, r-q*r2, u-q*u2, v-q*v2
	}
	//r = pgcd(a,b)
	// r = a*u+b*v
	if r != 1 {
		panic(fmt.Sprintf("gcd should be 1, got %d=gcd(%d, %d)", r, in.n, in.len))
	}
	u = (((u % in.len) + in.len) % in.len)
	return incN{u, in.len}
}

func TestShuffler(t *testing.T) {
	cc := []struct {
		op      Shuffler
		in, out int64
	}{
		{newStack{10}, 0, 9},
		{newStack{10}, 2, 7},
		{newStack{10}, 7, 2},

		{cutN{3, 10}, 0, 7},
		{cutN{3, 10}, 3, 0},
		{cutN{3, 10}, 9, 6},

		{cutN{10 - 4, 10}, 0, 4},
		{cutN{10 - 4, 10}, 9, 3},
		{cutN{10 - 4, 10}, 6, 0},

		{incN{3, 10}, 0, 0},
		{incN{3, 10}, 1, 3},
		{incN{3, 10}, 9, 7},
	}
	for _, c := range cc {
		t.Run(fmt.Sprintf("%#v", c.op), func(t *testing.T) {
			a := assert.New(t)
			out := c.op.NewPos(c.in)
			a.Equal(c.out, out)
			inv := c.op.Inverse()
			a.Equal(c.in, inv.NewPos(out))
		})
	}
}

func parseProcess(r io.Reader, l int) (process, error) {
	var shuf process

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "deal with increment ") {
			n, err := strconv.Atoi(strings.TrimPrefix(text, "deal with increment "))
			if err != nil {
				return shuf, err
			}
			shuf = append(shuf, incN{int64(n), int64(l)})
			continue
		}

		if strings.HasPrefix(text, "deal into new stack") {
			shuf = append(shuf, newStack{int64(l)})
			continue
		}

		if strings.HasPrefix(text, "cut ") {
			n, err := strconv.Atoi(strings.TrimPrefix(text, "cut "))
			if err != nil {
				return shuf, err
			}
			n = n % l
			if n < 0 {
				n += l
			}
			shuf = append(shuf, cutN{int64(n), int64(l)})
			continue
		}

		return shuf, fmt.Errorf("unsupported step: %s", text)
	}
	return shuf, scanner.Err()
}

func TestParse(t *testing.T) {
	cc := []struct {
		input    string
		len      int
		valueAt0 int64
		posOf0   int64
	}{
		{
			`deal with increment 7
deal into new stack
deal into new stack`, 3,
			0, 0,
		},
		{
			`cut 6
deal with increment 7
deal into new stack`, 3,
			3, 1,
		},
		{
			`deal with increment 7
deal with increment 9
cut -2`, 3,
			6, 2,
		},
		{
			`deal into new stack
cut -2
deal with increment 7
cut 8
cut -4
deal with increment 7
cut 3
deal with increment 9
deal with increment 3
cut -1`, 10,
			9, 7,
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			shuf, err := parseProcess(strings.NewReader(c.input), 10)
			a.NoError(err)
			a.Len(shuf, c.len)

			a.Equal(c.posOf0, shuf.NewPos(0))
			a.Equal(c.valueAt0, shuf.Inverse().NewPos(0))
		})
	}
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	shuf, err := parseProcess(f, 10007)
	a.NoError(err)

	a.Equal(int64(3589), shuf.NewPos(2019))
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	shuf, err := parseProcess(f, 119_315_717_514_047)
	a.NoError(err)

	inv := shuf.Inverse()

	p := int64(2020)
	N := 101_741_582_076_661
	cache := make(map[int64]bool)
	for i := 0; i < N; i++ {
		if cache[p] {
			fmt.Println(p)
			fmt.Println(shuf.NewPos(p))
			panic(p)
		}
		cache[p] = true
		if i%1_000_000 == 0 {
			fmt.Println(float64(100*i) / float64(N))
		}
		p = inv.NewPos(p)
	}
	a.Equal(3589, p)
}