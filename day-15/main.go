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
	grid := make([][]string, 0)
	moves := make([]string, 0)
	lineBreak := false
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			lineBreak = true
		} else {
			if !lineBreak {
				grid = append(grid, strings.Split(line, ""))
			} else {
				moves = append(moves, strings.Split(line, "")...)
			}
		}
	}

	var startTime = time.Now()
	result := part1(grid, moves)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	// startTime = time.Now()
	// result = part2(grid, moves)
	// fmt.Printf("Part 2: %d\n", result)
	// executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	// fmt.Printf("Completed Part 2 in %f seconds\n", executionTime)

}

type Coord struct {
	x int
	y int
}

func part1(grid [][]string, moves []string) int {

	robotPosition := findStart(grid)

	for _, move := range moves {
		var dir int
		if move == "^" || move == "v" {
			column := make([]string, len(grid))
			for y := 0; y < len(grid); y++ {
				column[y] = grid[y][robotPosition.x]
			}
			if move == "^" {
				dir = -1
			} else {
				dir = 1
			}
			robotPosition.y = push(&column, robotPosition.y, dir)
			//reassimilate the column back into the grid
			for y := 0; y < len(grid); y++ {
				grid[y][robotPosition.x] = column[y]
			}
		} else {
			row := grid[robotPosition.y]
			if move == "<" {
				dir = -1
			} else {
				dir = 1
			}
			robotPosition.x = push(&row, robotPosition.x, dir)
		}

	}

	return scoreGrid(grid)
}

func findStart(grid [][]string) Coord {

	for y, row := range grid {
		for x, ch := range row {
			if ch == "@" {
				return Coord{x: x, y: y}
			}
		}
	}

	return Coord{0, 0}
}

func push(arr *[]string, pos int, dir int) int {

	mark := -1
loop:
	for i := pos + dir; i >= 0 && i < len(*arr); i = i + dir {

		switch (*arr)[i] {
		case "#":
			// no room. break
			break loop
		case "O":
			// keep looking
		case ".":
			// space found. stop looping and shift
			mark = i
			break loop
		}
	}
	if mark > -1 {
		for i := mark; i != pos; i = i - dir {
			if (*arr)[i-dir] != "@" {
				(*arr)[i] = "O"
			} else {
				(*arr)[i] = "@"
			}
		}
		(*arr)[pos] = "."
		pos = pos + dir
	}

	return pos
}

func scoreGrid(grid [][]string) int {
	total := 0
	for y, row := range grid {
		for x, s := range row {
			fmt.Print(s)
			if s == "O" {
				total += ((100 * y) + x)
			}
		}
		fmt.Println()
	}
	return total
}
