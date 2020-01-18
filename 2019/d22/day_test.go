package d22

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type process []Shuffler

type generalisation struct {
	factor, shift int64
}

var identity = generalisation{1, 0}

func (g generalisation) compute(x int64, mod int64) int64 {
	x = (g.factor*x + g.shift) % mod
	if x < 0 {
		return x + mod
	}
	return x
}

func (g generalisation) power(x int64, mod int64, n int) int64 {
	factor := big.NewInt(g.factor)
	shift := big.NewInt(g.shift)
	y := big.NewInt(x)
	m := big.NewInt(mod)

	for n > 0 {
		if n%2 == 0 {
			// f*(f*x+s)+s
			shift.Add(new(big.Int).Mul(factor, shift), shift)
			shift.Mod(shift, m)

			factor.Mul(factor, factor)
			factor.Mod(factor, m)
			n = n / 2
		} else {
			y.Mod(new(big.Int).Add(new(big.Int).Mul(factor, y), shift), m)
			n--
		}
	}
	return y.Int64()
}

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

func (p process) Generalize(g generalisation) generalisation {
	for _, s := range p {
		g = s.Generalize(g)
	}
	return g
}

type Shuffler interface {
	NewPos(int64) int64
	Inverse() Shuffler
	Generalize(g generalisation) generalisation
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

func (ns newStack) Generalize(g generalisation) generalisation {
	g.factor = -g.factor
	g.shift = (ns.len - g.shift - 1) % ns.len
	return g
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
func (cn cutN) Generalize(g generalisation) generalisation {
	g.shift = (g.shift - cn.n) % cn.len
	return g
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

func (in incN) Generalize(g generalisation) generalisation {
	n := big.NewInt(in.n)
	l := big.NewInt(in.len)
	// g.factor = (g.factor * in.n) % in.len
	g.factor = new(big.Int).Mod(new(big.Int).Mul(big.NewInt(g.factor), n), l).Int64()
	// g.shift = (g.shift * in.n) % in.len
	g.shift = new(big.Int).Mod(new(big.Int).Mul(big.NewInt(g.shift), n), l).Int64()
	return g
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

			g := c.op.Generalize(identity)
			out = g.compute(c.in, 10)
			t.Log(g)
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
			a.Equal(c.posOf0, shuf.Generalize(identity).compute(0, 10))

			a.Equal(c.valueAt0, shuf.Inverse().NewPos(0))
			a.Equal(c.valueAt0, shuf.Inverse().Generalize(identity).compute(0, 10))
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

	g := shuf.Generalize(identity)

	t.Log(g)

	a.Equal(int64(3589), g.compute(2019, 10007))
	a.Equal(int64(3589), g.power(2019, 10007, 1))

	gi := shuf.Inverse().Generalize(identity)
	a.Equal(int64(2019), gi.power(3589, 10007, 1))
	a.Equal(int64(1), gi.factor*g.factor%10007)
	a.Equal(int64(0), (gi.factor*g.shift+gi.shift)%10007)
	a.Equal(int64(0), (g.factor*gi.shift+g.shift)%10007)

	a.Equal(int64(2019), gi.power(g.power(2019, 10007, 20), 10007, 20))
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	shuf, err := parseProcess(f, 119_315_717_514_047)
	a.NoError(err)

	g := shuf.Inverse().Generalize(identity)

	p := g.power(2020, 119_315_717_514_047, 101_741_582_076_661)
	a.Equal(int64(4893716342290), p)

	// gi := shuf.Generalize(identity)
	// t.Log(gi)

	// a.Equal(int64(1), gi.factor*g.factor%119_315_717_514_047)
	// a.Equal(int64(0), (gi.factor*g.shift+gi.shift)%119_315_717_514_047)
	// a.Equal(int64(0), (g.factor*gi.shift+g.shift)%119_315_717_514_047)

	// t.Log(gi.factor * g.factor % 119_315_717_514_047)
	// n := 1
	// a.Equal(2020, g.power(gi.power(2020, 119_315_717_514_047, n), 119_315_717_514_047, n))
}
