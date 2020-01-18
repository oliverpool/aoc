package d23

import (
	"fmt"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type packet struct {
	x, y int
}

func runComputer(intcodes map[int]int, addr int, send func(int, packet), recv <-chan packet, stalled *int32) error {
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
		wasStalled := false
		isStalled := false
		for {
			isStalled = false
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
				case <-time.After(time.Millisecond):
					select {
					case incoming = <-recv:
						buffer = append(buffer, incoming)
					case pInput <- -1:
						isStalled = true
					}
				}
			}
			if wasStalled != isStalled {
				if isStalled {
					atomic.AddInt32(stalled, 1)
				} else {
					atomic.AddInt32(stalled, -1)
				}
				wasStalled = isStalled
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
	stalled := new(int32)

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
		go runComputer(intcodes, i, send, recv[i], stalled)
	}

	res := <-final

	a.Equal(24106, res)
}

func TestSecond(t *testing.T) {
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
	stalled := new(int32)

	recv := make([]chan packet, 0, n)
	for i := 0; i < n; i++ {
		recv = append(recv, make(chan packet, 10))
	}

	// var nat packet
	var naty = new(int32)
	var natx = new(int32)
	send := func(addr int, p packet) {
		if addr == 255 {
			atomic.StoreInt32(natx, int32(p.x))
			atomic.StoreInt32(naty, int32(p.y))
			return
		}
		recv[addr] <- p
	}

	// idle management
	go func() {
		lasty := -1
		for range time.Tick(time.Millisecond) {
			s := atomic.LoadInt32(stalled)
			fmt.Println(s)
			if s == 50 {
				time.Sleep(time.Millisecond)
				s = atomic.LoadInt32(stalled)
				if s != 50 {
					fmt.Println("false alert", s)
					continue
				}
				nx := int(atomic.LoadInt32(natx))
				ny := int(atomic.LoadInt32(naty))
				if lasty == ny {
					final <- lasty
				}
				send(0, packet{nx, ny})
				lasty = ny

				for atomic.LoadInt32(stalled) == 50 {
					time.Sleep(time.Millisecond)
				}
			}
		}
	}()

	for i := 0; i < n; i++ {
		go runComputer(intcodes, i, send, recv[i], stalled)
	}

	res := <-final

	a.Equal(17895, res)
}
