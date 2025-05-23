package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Maze [][]string

func (maze Maze) At(c Coord) (string, error) {
	if !maze.IsInBounds(c) {
		return "", errors.New("Coordinate out of bounds")
	}
	return maze[c.row][c.col], nil
}

func (maze Maze) IsWall(c Coord) bool {
	if !maze.IsInBounds(c) {
		return false
	}
	return maze[c.row][c.col] == "#"
}

func (maze Maze) IsSpace(c Coord) bool {
	if !maze.IsInBounds(c) {
		return false
	}
	return maze[c.row][c.col] != "#"
}

func (maze Maze) IsInBounds(c Coord) bool {
	return c.row >= 0 && c.row < len(maze) && c.col >= 0 && c.col < len(maze[0])
}

type Coord struct {
	row int
	col int
}

var directions = []Coord{
	{0, 1},
	{1, 0},
	{0, -1},
	{-1, 0},
}

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
	maze := make([][]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, strings.Split(line, ""))
	}

	var startTime = time.Now()
	result := part1(maze)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	startTime = time.Now()
	result = part2(maze)
	fmt.Printf("Part 1: %d\n", result)
	executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)
}

func part1(maze Maze) int {
	var start, target Coord
	for r, row := range maze {
		for c, str := range row {
			if str == "S" {
				start = Coord{row: r, col: c}
			} else if str == "E" {
				target = Coord{row: r, col: c}
			}
		}
	}

	// The path is pre-ordained. No forks. So walk the path and assign each coordinate in the
	// path it's step value
	path := []Coord{start}
	distances := make(map[Coord]int)
	distances[start] = 0
	loc := start
	steps := 0
	for loc != target {
		for _, d := range directions {
			c := move(loc, d, 1)
			if _, found := distances[c]; !found && maze[c.row][c.col] != "#" {
				loc = c
				steps++
				distances[loc] = steps
				path = append(path, loc)
				break
			}
		}
	}

	// For each step in the path
	//  is there a 1 or 2 space wall nearby?
	//  if we removed the wall, how many steps would be saved
	//    find difference between the two newly connected spaces,
	//    and account for the wall size
	total := 0
	saves := make(map[int]int)
	for _, pathStep := range path {

		for _, dir := range directions {
			wallSize := wallSize(maze, pathStep, dir)
			if wallSize == 1 || wallSize == 2 {
				dest := move(pathStep, dir, wallSize+1)
				if destSteps, found := distances[dest]; found {
					saved := destSteps - distances[pathStep] - 1 - wallSize
					if saved >= 100 {
						saves[saved] = saves[saved] + 1
						total += 1
					}
				}
			}
		}
	}

	return total
}

func part2(maze Maze) int {
	var start, target Coord
	for r, row := range maze {
		for c, str := range row {
			if str == "S" {
				start = Coord{row: r, col: c}
			} else if str == "E" {
				target = Coord{row: r, col: c}
			}
		}
	}

	// The path is pre-ordained. No forks. So walk the path and assign each coordinate in the
	// path it's step value
	path := []Coord{start}
	distances := make(map[Coord]int)
	distances[start] = 0
	loc := start
	steps := 0
	for loc != target {
		for _, d := range directions {
			c := move(loc, d, 1)
			if _, found := distances[c]; !found && maze[c.row][c.col] != "#" {
				loc = c
				steps++
				distances[loc] = steps
				path = append(path, loc)
				break
			}
		}
	}

	total := 0
	for _, loc1 := range path {
		spaces := spacesInRange(maze, loc1, 20)
		for _, loc2 := range spaces {
			if loc1.row == 4 && loc1.col == 4 && loc2.row == 5 && loc2.col == 5 {
				fmt.Println("Here")
			}
			dist1 := distances[loc1]
			dist2 := distances[loc2]
			if dist1 < dist2 {
				cheatSteps := abs(loc1.row-loc2.row) + abs(loc1.col-loc2.col)
				saved := dist2 - dist1 - cheatSteps
				if saved >= 100 {
					total += 1
				}
			}
		}
	}

	return total
}

func abs(n int) int {
	if n < 0 {
		n = n * -1
	}
	return n
}

func spacesInRange(maze Maze, loc Coord, rng int) []Coord {
	spacesMap := make(map[Coord]bool)

	submitSpace := func(space Coord) {
		if space != loc && !spacesMap[space] && maze.IsSpace(space) {
			d := abs(space.row-loc.row) + abs(space.col-loc.col)
			if d > rng {
				fmt.Printf("%v is actually %d away\n", space, d)
			}
			spacesMap[space] = true
		}
	}

	for colOffset := 0; colOffset <= rng; colOffset++ {
		for rowOffset := 0; rowOffset <= rng-colOffset; rowOffset++ {
			space := Coord{row: loc.row + rowOffset, col: loc.col + colOffset}
			submitSpace(space)

			space = Coord{row: loc.row - rowOffset, col: loc.col + colOffset}
			submitSpace(space)

			space = Coord{row: loc.row - rowOffset, col: loc.col - colOffset}
			submitSpace(space)

			space = Coord{row: loc.row + rowOffset, col: loc.col - colOffset}
			submitSpace(space)
		}

	}
	spaces := make([]Coord, 0)
	for space := range spacesMap {
		spaces = append(spaces, space)
	}
	return spaces

}

func wallSize(maze Maze, loc, dir Coord) int {
	d := loc
	wallSize := 0
	for true {
		d = move(d, dir, 1)
		if maze.IsWall(d) {
			wallSize++
		} else {
			return wallSize
		}
	}

	return wallSize
}

func move(loc, dir Coord, steps int) Coord {
	dest := loc
	for i := 0; i < steps; i++ {
		dest = Coord{row: dest.row + dir.row, col: dest.col + dir.col}
	}
	return dest
}
