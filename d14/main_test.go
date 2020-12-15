package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X
mem[8] = 11
mem[7] = 101
mem[8] = 0`

type state struct {
	mem  map[uint64]uint64
	mask mask
}

func (s state) String() string {
	return fmt.Sprintf("%+v", s.mem)
}

type instruction interface {
	execute(*state)
	executev2(*state)
}

func newMem(addr, value string) mem {
	a, err := strconv.ParseUint(addr, 10, 64)
	if err != nil {
		panic(err)
	}
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return mem{
		address: a,
		value:   v,
	}
}

type mem struct {
	address uint64
	value   uint64
}

func (m mem) execute(s *state) {
	previous := m.value
	s.mem[m.address] = (previous&s.mask.and | s.mask.or)
}
func (m mem) executev2(s *state) {
	addr := (m.address | s.mask.and) ^ (s.mask.or ^ s.mask.and)
	for _, f := range s.mask.floatings() {
		s.mem[addr|f] = m.value
		// fmt.Printf("% 08b\n", addr)
		// fmt.Printf("% 08b => %d\n", addr|f, m.value)
		// for _, f := range ff {
		// 	fmt.Printf("% 08b\n", f)
		// }
	}
}

func newMask(s string) mask {
	sAnd := strings.ReplaceAll(s, "X", "1")
	mAnd, err := strconv.ParseUint(sAnd, 2, 64)
	if err != nil {
		panic(err)
	}

	sOr := strings.ReplaceAll(s, "X", "0")
	mOr, err := strconv.ParseUint(sOr, 2, 64)
	if err != nil {
		panic(err)
	}
	return mask{
		or:  mOr,
		and: mAnd,
	}
}

type mask struct {
	or  uint64
	and uint64
}

func (m mask) execute(s *state) {
	s.mask = m
}

func (m mask) executev2(s *state) {
	s.mask = m
}

func enumerate(xor uint64) []uint64 {
	if xor == 0 {
		return nil
	}
	if xor == 1 {
		return []uint64{0, 1}
	}
	if xor%2 == 0 {
		got := enumerate(xor / 2)
		ff := make([]uint64, 0, len(got))
		for _, g := range got {
			ff = append(ff, 2*g)
		}
		return ff
	}
	got := enumerate((xor - 1) / 2)
	ff := make([]uint64, 0, 2*len(got))
	for _, g := range got {
		ff = append(ff, 2*g, 2*g+1)
	}
	return ff
}

func (m mask) floatings() []uint64 {
	xor := m.or ^ m.and
	ff := enumerate(xor)
	// fmt.Printf("=> % 08b\n", xor)
	// for _, f := range ff {
	// 	fmt.Printf("% 08b\n", f)
	// }
	return ff
}

func parseInput(s string) (inst []instruction) {
	for _, l := range strings.Split(strings.TrimSpace(s), "\n") {
		parts := strings.Split(l, " ")
		if parts[0] == "mask" {
			inst = append(inst, newMask(parts[2]))
		} else if l[:3] == "mem" {
			addr := parts[0][len("mem[") : len(parts[0])-1]
			inst = append(inst, newMem(addr, parts[2]))
		} else {
			panic(l)
		}
	}
	return inst
}

func runAndSum(inst []instruction) uint64 {
	s := state{
		mem: make(map[uint64]uint64),
	}
	for _, i := range inst {
		i.execute(&s)
	}
	var sum uint64
	for _, v := range s.mem {
		sum += v
	}
	return sum
}

func TestFirst(t *testing.T) {
	prog := parseInput(example)
	s := runAndSum(prog)
	require.Equal(t, uint64(165), s)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	prog = parseInput(string(input))
	s = runAndSum(prog)
	require.Equal(t, uint64(14954914379452), s)
}

func runAndSumVersion2(inst []instruction) uint64 {
	s := state{
		mem: make(map[uint64]uint64),
	}
	for _, i := range inst {
		i.executev2(&s)
	}
	var sum uint64
	for _, v := range s.mem {
		sum += v
	}
	return sum
}

func TestSecond(t *testing.T) {
	prog := parseInput(`mask = 000000000000000000000000000000X1001X
mem[42] = 100
mask = 00000000000000000000000000000000X0XX
mem[26] = 1`)
	s := runAndSumVersion2(prog)
	t.Log(s)
	require.Equal(t, uint64(208), s)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	prog = parseInput(string(input))
	s = runAndSumVersion2(prog)
	t.Log(s)
	// 135703051814: to low
	require.Equal(t, uint64(3415488160714), s)
}
