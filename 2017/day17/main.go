package main

import (
	"fmt"
)

func main() {
	pos := computePositions(9, 3)

	fmt.Println(pos)
	fmt.Println(values(pos))
	fmt.Println(valueAfterZero(pos))

	test := computePositions(2017, 316)
	fmt.Println(valueAfterLast(test))

	test2 := computePositions(50000000, 316)
	fmt.Println(valueAfterZero(test2))
}

func computePositions(length, step int) []int {
	positions := make([]int, length+1)
	for i := 1; i <= length; i++ {
		positions[i] = ((positions[i-1] + step) % i) + 1
	}
	return positions
}

func valueAfterZero(positions []int) int {
	for index := len(positions) - 1; index > 0; index-- {
		val := positions[index]
		if val == 1 {
			return index
		}
	}
	panic("not found")
}

func valueAfterLast(positions []int) int {
	pos := positions[len(positions)-1]
	for index := len(positions) - 2; index > 0; index-- {
		val := positions[index]
		if val == pos {
			return index
		} else if val < pos {
			pos--
		}
	}
	panic("not found")
}

func values(positions []int) (values []int) {
	for val, pos := range positions {
		values = append(values[0:pos], append([]int{val}, values[pos:]...)...)
	}
	return values
}
