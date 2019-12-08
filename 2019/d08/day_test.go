package d08

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findLayer(img []byte, size int) int {
	count := make(map[byte]int)
	minCount := -1
	minValue := 0
	for i, b := range img {
		count[b]++
		if (i+1)%size == 0 {
			if minCount == -1 || count['0'] < minCount {
				minCount = count['0']
				minValue = count['1'] * count['2']
			}
			count = make(map[byte]int)
		}
	}
	return minValue
}

func TestLayer(t *testing.T) {
	testCases := []struct {
		input  string
		size   int
		output int
	}{
		{
			input:  "123456789012",
			size:   3 * 2,
			output: 1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			a := assert.New(t)
			n := findLayer([]byte(tC.input), tC.size)
			a.Equal(tC.output, n)
		})
	}
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	img, err := ioutil.ReadAll(f)
	a.NoError(err)

	out := findLayer(img, 25*6)
	a.Equal(2064, out)
}

func stackedLayers(img []byte, size int) []byte {
	result := append([]byte(nil), img[0:size]...)

	for i, b := range img[size:] {
		ii := i % size
		if result[ii] != '2' {
			continue
		}
		result[ii] = b
	}

	return result
}

func TestStackedLayer(t *testing.T) {
	testCases := []struct {
		input  string
		size   int
		output string
	}{
		{
			input:  "0222112222120000",
			size:   2 * 2,
			output: "0110",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			a := assert.New(t)
			n := stackedLayers([]byte(tC.input), tC.size)
			a.Equal(tC.output, string(n))
		})
	}
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	img, err := ioutil.ReadAll(f)
	a.NoError(err)

	w := 25
	out := string(stackedLayers(img, w*6))

	i := 0
	expected := []string{
		"8  8  88  8  8 8888  88  ",
		"8 8  8  8 8  8    8 8  8 ",
		"88   8  8 8  8   8  8  8 ",
		"8 8  8888 8  8  8   8888 ",
		"8 8  8  8 8  8 8    8  8 ",
		"8  8 8  8  88  8888 8  8 ",
	}
	for i*w < len(out) {
		line := out[i*w : i*w+w]
		line = strings.ReplaceAll(line, "1", "8")
		line = strings.ReplaceAll(line, "0", " ")
		a.Equal(expected[i], line)
		i++
	}
}
