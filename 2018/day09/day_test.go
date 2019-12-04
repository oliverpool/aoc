package day_test

import (
	"fmt"
	"testing"
)

func TestOne(t *testing.T) {
	players := make([]int, 470)
	lastMarble := 72170

	board := make([]int, 1, lastMarble)

	insertMarble := func(marble, at int) {
		pos := at
		var value int
		for pos < len(board) {
			board[pos], value = value, board[pos]
			pos += 1
		}
		board = append(board, value)
		board[at] = marble
	}

	index := 1
	for marble := 1; marble <= lastMarble; marble++ {
		if marble%23 == 0 {
			currentPlayer := marble % len(players)
			players[currentPlayer] += marble
			index = (index + len(board) - 7) % len(board)
			players[currentPlayer] += board[index]
			board = append(board[:index], board[index+1:]...)
			continue
		}
		index = ((index + 1) % len(board)) + 1
		insertMarble(marble, index)
	}
	max := 0
	for _, m := range players {
		if m > max {
			max = m
		}
	}
	t.Fatal(max)
}

func TestTwo(t *testing.T) {
	// defer profile.Start().Stop()
	players := make([]int, 470)
	lastMarble := 72170 * 100

	indexes := []int{0}

	index := 1
	actualMarbles := 1
	for marble := 1; marble <= lastMarble; marble++ {
		if marble%100000 == 0 {
			fmt.Println(lastMarble - marble)
		}
		// t.Log(indexes)
		// t.Log(buildBoard(indexes))
		// t.Log(marble, indexes[len(indexes)-1])

		if marble%23 == 0 {
			currentPlayer := marble % len(players)
			players[currentPlayer] += marble
			index = (index + actualMarbles - 7) % actualMarbles
			// t.Log(indexes)

			indexes = append(indexes, -1)
			// find the value @ index
			currentIndex := index
			for i := len(indexes) - 1; i >= 0; i-- {
				if indexes[i] == currentIndex {
					players[currentPlayer] += i
					indexes[i] = -1
					// t.Log(i)
					break
				} else if indexes[i] >= 0 && indexes[i] < currentIndex {
					currentIndex -= 1
				} else {
					indexes[i] -= 1
				}
			}
			actualMarbles -= 1
			continue
		}
		index = ((index + 1) % actualMarbles) + 1
		indexes = append(indexes, index)
		actualMarbles += 1
		// t.Log(indexes)
	}
	max := 0
	for _, m := range players {
		if m > max {
			max = m
		}
	}
	t.Fatal(max)
}

func buildBoard(indexes []int) []int {
	var board []int

	insertMarble := func(marble, at int) {
		pos := at
		var value int
		for pos < len(board) {
			board[pos], value = value, board[pos]
			pos += 1
		}
		board = append(board, value)
		board[at] = marble
	}

	for i, v := range indexes {
		if v < 0 {
			continue
		}
		insertMarble(i, v)
	}
	return board
}
