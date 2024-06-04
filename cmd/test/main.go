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

func checkWinOriginal(board [size][size]int, player int) bool {
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

func checkWinOptimized(board [size][size]int, player int) bool {
	var hor, ver, diag1, diag2 [size][size]int

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if board[i][j] == player {
				hor[i][j] = 1
				ver[i][j] = 1
				diag1[i][j] = 1
				diag2[i][j] = 1

				if j > 0 {
					hor[i][j] += hor[i][j-1]
				}
				if i > 0 {
					ver[i][j] += ver[i-1][j]
				}
				if i > 0 && j > 0 {
					diag1[i][j] += diag1[i-1][j-1]
				}
				if i > 0 && j < size-1 {
					diag2[i][j] += diag2[i-1][j+1]
				}

				if hor[i][j] >= winCount || ver[i][j] >= winCount || diag1[i][j] >= winCount || diag2[i][j] >= winCount {
					return true
				}
			}
		}
	}

	return false
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
		checkWinOriginal(board, cross)
	}
	duration := time.Since(start)
	fmt.Printf("Оригинальный алгоритм: Время выполнения 1000 раз: %d наносекунд\n", duration.Nanoseconds())

	start = time.Now()
	for i := 0; i < 1000; i++ {
		checkWinOptimized(board, cross)
	}
	duration = time.Since(start)
	fmt.Printf("Оптимизированный алгоритм: Время выполнения 1000 раз: %d наносекунд\n", duration.Nanoseconds())
}
