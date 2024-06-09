package main

import (
	"fmt"
	"time"
)

const (
	size     = 1000
	empty    = 0
	cross    = 1
	nought   = 2
	winCount = 5
)

var directions = [][2]int{
	{1, 0},
	{0, 1},
	{1, 1},
	{1, -1},
}

func checkWin(board [size][size]int, player int) bool {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if board[i][j] == player {
				for _, dir := range directions {
					if checkDirection(board, player, i, j, dir[0], dir[1]) {
						return true
					}
				}
			}
		}
	}
	return false
}

func checkDirection(board [size][size]int, player, x, y, dx, dy int) bool {
	count := 0
	for k := 0; k < winCount; k++ {
		nx, ny := x+dx*k, y+dy*k
		if nx < 0 || nx >= size || ny < 0 || ny >= size || board[nx][ny] != player {
			return false
		}
		count++
	}
	return count == winCount
}

func main() {
	var board [size][size]int
	board[size-1][size-1] = cross
	board[size-2][size-2] = cross
	board[size-3][size-3] = cross
	board[size-4][size-4] = cross
	board[size-5][size-5] = cross

	start := time.Now()
	for i := 0; i < 1000; i++ {
		checkWin(board, cross)
	}
	duration := time.Since(start)
	fmt.Printf("Оригинальный алгоритм: Время выполнения 1000 раз: %d наносекунд\n", duration.Nanoseconds())
}
