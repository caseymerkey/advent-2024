package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	inputFile := "sample.txt"
	if len(os.Args) > 1 && len(os.Args[1]) > 0 {
		inputFile = os.Args[1]
	}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	puzzle := make([][]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")
		puzzle = append(puzzle, row)
	}

	var startTime = time.Now()
	result := part1(puzzle)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	startTime = time.Now()
	result = part2(puzzle)
	fmt.Printf("Part 2: %d\n", result)
	executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 2 in %f seconds\n", executionTime)
}

const (
	N = iota
	E
	S
	W
)

func part1(puzzleBase [][]string) int {

	traversedCount := 0
	r, c := 0, 0

	direction := N
	edge := len(puzzleBase) - 1 // thankfully, it's a square

	puzzle := make([][]string, 0)

	for row, cells := range puzzleBase {
		puzzleRow := make([]string, 0)
		for col, cell := range cells {
			if cell == "^" {
				r = row
				c = col
			}
			puzzleRow = append(puzzleRow, cell)
		}
		puzzle = append(puzzle, puzzleRow)
	}

	exiting := false
	for !exiting {
		pathIsClear := false
		for !pathIsClear {
			r1 := r
			c1 := c
			switch direction {
			case N:
				r1--
			case S:
				r1++
			case E:
				c1++
			case W:
				c1--
			}

			if puzzle[r1][c1] == "#" {
				direction++
				if direction > W {
					direction = 0
				}
			} else {
				pathIsClear = true
				r = r1
				c = c1
				if puzzle[r][c] != "X" {
					puzzle[r][c] = "X"
					traversedCount++
				}
			}
		}

		if (r == 0 && direction == N) || (r == edge && direction == S) || (c == 0 && direction == W) || (c == edge && direction == E) {
			exiting = true
		}
	}

	return traversedCount
}

func part2(puzzle [][]string) int {
	return 0
}
