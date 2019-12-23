package d23

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type packet struct {
	x, y int
}

func runComputer(intcodes map[int]int, addr int, send func(int, packet), recv <-chan packet) error {
	copy := make(map[int]int)
	for i, v := range intcodes {
		copy[i] = v
	}

	pOutput := make(chan int, 3)
	go func() {
		var p packet
		for {
			addr := <-pOutput
			p.x = <-pOutput
			p.y = <-pOutput
			send(addr, p)
		}
	}()

	pInput := make(chan int)
	go func() {
		pInput <- addr
		var buffer []packet
		var current packet
		var incoming packet
		hasOne := false
		for {
			if !hasOne && len(buffer) > 0 {
				current, buffer = buffer[0], buffer[1:]
				hasOne = true
			}
			if hasOne {
				select {
				case incoming = <-recv:
					buffer = append(buffer, incoming)
				case pInput <- current.x:
					pInput <- current.y
					hasOne = false
				}
			} else {
				select {
				case incoming = <-recv:
					buffer = append(buffer, incoming)
				case pInput <- -1:
				}
			}
		}
	}()

	return runProgram(copy, pInput, pOutput)
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

	final := make(chan int, 1)

	n := 50

	recv := make([]chan packet, 0, n)
	for i := 0; i < n; i++ {
		recv = append(recv, make(chan packet, 10))
	}

	send := func(addr int, p packet) {
		if addr == 255 {
			final <- p.y
			return
		}
		recv[addr] <- p
	}

	for i := 0; i < n; i++ {
		go runComputer(intcodes, i, send, recv[i])
	}

	res := <-final

	a.Equal(24106, res)
}
