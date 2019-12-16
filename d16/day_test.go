package d16

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func fftSlow(ins string, phases, skip int) string {
	in := make([]int, 0, len(ins))
	for _, b := range ins {
		in = append(in, int(b-'0'))
	}

	in2 := make([]int, len(ins))
	for phases > 0 {
		// if phases%10 == 0 {
		fmt.Println(phases)
		// }
		for i := range in {
			in2[i] = fftOneFast(in, skip+i, skip)
		}
		copy(in, in2)
		phases--
	}

	out := make([]byte, 0, len(in))
	for _, i := range in {
		out = append(out, '0'+byte(i))
	}
	return string(out)
}

func fft(ins string, phases, skip int) string {
	in := make([]int, 0, len(ins))
	for _, b := range ins {
		in = append(in, int(b-'0'))
	}

	sums := make([]int, len(in)+1)
	for phases > 0 {
		if phases%10 == 0 {
			// fmt.Println(phases)
		}

		sum := 0
		for i, v := range in {
			sum += v
			sums[i+1] = sum
		}

		for i := range in {
			in[i] = fftOneFaster(sums, skip+i, skip)
		}
		phases--
	}

	out := make([]byte, 0, 8)
	for _, i := range in[:8] {
		out = append(out, '0'+byte(i))
	}
	return string(out)
}

func getFactor(i, s int) int {
	t := ((i + 1) / (s + 1)) % 4
	if t > 1 {
		return 2 - t
	}
	return t
}

func fftOne(in []int, shift int) int {
	out := 0
	for i, v := range in {
		s := ((i + 1) / (shift + 1)) % 4
		if s > 1 {
			s = 2 - s
		}
		out += s * v
	}
	if out < 0 {
		out = -out
	}
	return out % 10
}

func fftOneFaster(sums []int, shift, skip int) int {
	out := 0
	li := len(sums)
	for first := shift - skip; first < li; first += 4 * (shift + 1) {

		// add ones
		end := first + shift + 1
		if end >= li {
			end = li - 1
		}
		out += sums[end] - sums[first]

		// substract ones
		f := first + 2*(shift+1)
		if f >= li {
			continue
		}
		end = f + shift + 1
		if end >= li {
			end = li - 1
		}
		out -= sums[end] - sums[f]
	}

	if out < 0 {
		out = -out
	}
	return out % 10
}

func fftOneFast(in []int, shift, skip int) int {
	out := 0
	for first := shift - skip; first < len(in); first += 4 * (shift + 1) {
		// add ones
		for i := first; i <= first+shift && i < len(in); i++ {
			out += in[i]
		}

		// add -ones
		f := first + 2*(shift+1)
		for i := f; i <= f+shift && i < len(in); i++ {
			out -= in[i]
		}
	}

	if out < 0 {
		out = -out
	}
	return out % 10
}

func fftTenThousand(ins string, phases int) string {
	skip, _ := strconv.Atoi(ins[:7])
	ins = strings.Repeat(ins, 10000)

	// a number final value only depend on following numbers
	ins = ins[skip:]

	return fft(ins, phases, skip)
}

func TestFFTOne(t *testing.T) {
	cc := []struct {
		in    []int
		shift int
		out   int
	}{
		{
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			0,
			4,
		},
		{
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			1,
			8,
		},
		{
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			2,
			2,
		},
		{
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			3,
			2,
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			out := fftOne(c.in, c.shift)
			a.Equal(c.out, out)

			out = fftOneFast(c.in, c.shift, 0)
			a.Equal(c.out, out, "fast")
		})
	}
}

func TestNPhases(t *testing.T) {
	cc := []struct {
		input   string
		phases  int
		output8 string
	}{
		{
			"12345678",
			1,
			"48226158",
		},
		{
			"12345678",
			4,
			"01029498",
		},
		{
			"80871224585914546619083218645595",
			100,
			"24176176",
		},
		{
			"19617804207202209144916044189917",
			100,
			"73745418",
		},
		{
			"69317163492948606335995924319873",
			100,
			"52432133",
		},
	}
	for _, c := range cc {
		t.Run(c.input, func(t *testing.T) {
			a := assert.New(t)

			output := fft(c.input, c.phases, 0)
			a.Equal(c.output8, output)
		})
	}
}

func TestFirst(t *testing.T) {
	a := assert.New(t)

	in, err := ioutil.ReadFile("./input")
	a.NoError(err)
	ins := strings.TrimSpace(string(in))

	out := fft(ins, 100, 0)

	a.Equal("12541048", out)
}

func TestThousand(t *testing.T) {
	cc := []struct {
		input   string
		phases  int
		output8 string
	}{
		{
			"03036732577212944063491565474664",
			100,
			"84462026",
		},
		{
			"02935109699940807407585447034323",
			100,
			"78725270",
		},
		{
			"03081770884921959731165446850517",
			100,
			"53553731",
		},
	}
	for _, c := range cc {
		t.Run(c.input, func(t *testing.T) {
			a := assert.New(t)

			output := fftTenThousand(c.input, c.phases)
			a.Equal(c.output8, output)
		})
	}
}

func TestSecond(t *testing.T) {
	a := assert.New(t)

	in, err := ioutil.ReadFile("./input")
	a.NoError(err)
	ins := strings.TrimSpace(string(in))

	out := fftTenThousand(ins, 100)

	a.Equal("62858988", out)
}
