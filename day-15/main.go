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
	gridBase := make([][]string, 0)
	moves := make([]string, 0)
	lineBreak := false
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			lineBreak = true
		} else {
			if !lineBreak {
				gridBase = append(gridBase, strings.Split(line, ""))
			} else {
				moves = append(moves, strings.Split(line, "")...)
			}
		}
	}
	grid := make([][]string, len(gridBase))
	for y, row := range gridBase {
		grid[y] = make([]string, len(row))
		for x, c := range row {
			grid[y][x] = c
		}
	}

	var startTime = time.Now()
	result := part1(grid, moves)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	startTime = time.Now()
	result = part2(gridBase, moves)
	fmt.Printf("Part 2: %d\n", result)
	executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 2 in %f seconds\n", executionTime)

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

func part2(grid [][]string, moves []string) int {
	newGrid := make([][]string, len(grid))
	for y, row := range grid {
		newRow := make([]string, len(row)*2)
		for x, cell := range row {
			if cell == "@" {
				newRow[x*2] = "@"
				newRow[x*2+1] = "."
			} else if cell == "O" {
				newRow[x*2] = "["
				newRow[x*2+1] = "]"
			} else {
				newRow[x*2] = cell
				newRow[x*2+1] = cell
			}
		}
		newGrid[y] = newRow
	}
	grid = newGrid
	robotPosition := findStart(grid)

	for _, move := range moves {
		var dir int
		if move == "^" || move == "v" {

			if move == "^" {
				dir = -1
			} else {
				dir = 1
			}
			target := grid[robotPosition.y+dir][robotPosition.x]
			switch target {
			case ".":
				grid[robotPosition.y+dir][robotPosition.x] = "@"
				grid[robotPosition.y][robotPosition.x] = "."
				robotPosition.y = robotPosition.y + dir
			case "#":
				// blocked. No movement
			case "[":
				if verticalPushAllowed(&grid, Coord{x: robotPosition.x, y: robotPosition.y + dir}, dir) {
					pushVertical(&grid, Coord{x: robotPosition.x, y: robotPosition.y + dir}, dir)
					grid[robotPosition.y+dir][robotPosition.x] = "@"
					grid[robotPosition.y][robotPosition.x] = "."
					robotPosition.y = robotPosition.y + dir
				}

			case "]":
				if verticalPushAllowed(&grid, Coord{x: robotPosition.x - 1, y: robotPosition.y + dir}, dir) {
					pushVertical(&grid, Coord{x: robotPosition.x - 1, y: robotPosition.y + dir}, dir)
					grid[robotPosition.y+dir][robotPosition.x] = "@"
					grid[robotPosition.y][robotPosition.x] = "."
					robotPosition.y = robotPosition.y + dir
				}
			}

		} else {
			if move == "<" {
				dir = -1
			} else {
				dir = 1
			}
			target := grid[robotPosition.y][robotPosition.x+dir]
			if target == "#" {
				// blocked
			} else if target == "." {
				grid[robotPosition.y][robotPosition.x+dir] = "@"
				grid[robotPosition.y][robotPosition.x] = "."
				robotPosition.x = robotPosition.x + dir
			} else {
				if pushHorizontal(&grid, Coord{y: robotPosition.y, x: robotPosition.x + dir}, dir) {
					grid[robotPosition.y][robotPosition.x+dir] = "@"
					grid[robotPosition.y][robotPosition.x] = "."
					robotPosition.x = robotPosition.x + dir
				}
			}
		}

	}

	return scoreGrid(grid)
}

func pushHorizontal(grid *[][]string, boxLoc Coord, dir int) bool {
	pushCompleted := false

	if (*grid)[boxLoc.y][boxLoc.x+(dir*2)] == "#" {
		// blocked
		pushCompleted = false
	} else if (*grid)[boxLoc.y][boxLoc.x+(dir*2)] == "." {
		pushCompleted = true
	} else {
		pushCompleted = pushHorizontal(grid, Coord{x: boxLoc.x + (dir * 2), y: boxLoc.y}, dir)
	}
	if pushCompleted {
		(*grid)[boxLoc.y][boxLoc.x+(dir*2)] = (*grid)[boxLoc.y][boxLoc.x+dir]
		(*grid)[boxLoc.y][boxLoc.x+dir] = (*grid)[boxLoc.y][boxLoc.x]
		(*grid)[boxLoc.y][boxLoc.x] = "."
	}

	return pushCompleted
}

func pushVertical(grid *[][]string, boxLoc Coord, dir int) {

	c1 := Coord{y: boxLoc.y + dir, x: boxLoc.x}
	c2 := Coord{y: boxLoc.y + dir, x: boxLoc.x + 1}

	if (*grid)[c1.y][c1.x] == "[" {
		pushVertical(grid, c1, dir)
	} else {

		if (*grid)[c1.y][c1.x] == "]" {
			pushVertical(grid, Coord{x: c1.x - 1, y: c1.y}, dir)
		}
		if (*grid)[c2.y][c2.x] == "[" {
			pushVertical(grid, Coord{x: c2.x, y: c2.y}, dir)
		}

	}

	(*grid)[c1.y][c1.x] = "["
	(*grid)[c2.y][c2.x] = "]"
	(*grid)[boxLoc.y][boxLoc.x] = "."
	(*grid)[boxLoc.y][boxLoc.x+1] = "."

}

func verticalPushAllowed(grid *[][]string, boxLoc Coord, dir int) bool {
	allowed := false

	c1 := Coord{y: boxLoc.y + dir, x: boxLoc.x}
	c2 := Coord{y: boxLoc.y + dir, x: boxLoc.x + 1}

	if (*grid)[c1.y][c1.x] == "#" || (*grid)[c2.y][c2.x] == "#" {
		allowed = false
	} else if (*grid)[c1.y][c1.x] == "." && (*grid)[c2.y][c2.x] == "." {
		allowed = true
	} else {
		if (*grid)[c1.y][c1.x] == "[" {
			allowed = verticalPushAllowed(grid, c1, dir)
		} else {
			tmp := true
			if (*grid)[c1.y][c1.x] == "]" {
				tmp = tmp && verticalPushAllowed(grid, Coord{x: c1.x - 1, y: c1.y}, dir)
			}
			if (*grid)[c2.y][c2.x] == "[" {
				tmp = tmp && verticalPushAllowed(grid, Coord{x: c2.x, y: c2.y}, dir)
			}
			allowed = tmp
		}
	}

	return allowed
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
			if s == "O" || s == "[" {
				total += ((100 * y) + x)
			}
		}
		fmt.Println()
	}
	return total
}
