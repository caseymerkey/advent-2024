package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var directions = []Coord{
	{-1, 0},
	{0, -1},
	{1, 0},
	{0, 1},
}

const pathEndVal = 9

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
	puzzle := make([][]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		numStrings := strings.Split(line, "")
		rowNums := make([]int, 0)
		for _, s := range numStrings {
			n, _ := strconv.Atoi(s)
			rowNums = append(rowNums, n)
		}
		puzzle = append(puzzle, rowNums)
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

type Coord struct {
	row int
	col int
}

func part1(puzzle [][]int) int {
	total := 0

	for r := 0; r < len(puzzle); r++ {
		for c := 0; c < len(puzzle[0]); c++ {
			if puzzle[r][c] == 0 {
				// start here
				coord := Coord{row: r, col: c}
				foundEnds := make([]Coord, 0)
				ends := findTrailEnds(puzzle, coord, &foundEnds)
				total += ends
			}
		}
	}

	return total
}

func findTrailEnds(puzzle [][]int, coord Coord, foundEnds *[]Coord) int {

	target := puzzle[coord.row][coord.col] + 1
	endCount := 0
	for _, dir := range directions {
		testCoord := Coord{row: coord.row + dir.row, col: coord.col + dir.col}
		if testCoord.row >= 0 && testCoord.row < len(puzzle) && testCoord.col >= 0 && testCoord.col < len(puzzle[0]) && puzzle[testCoord.row][testCoord.col] == target {
			if target == pathEndVal {
				if !slices.Contains(*foundEnds, testCoord) {
					*foundEnds = append(*foundEnds, testCoord)
					endCount++
				}
			} else {
				endsFromHere := findTrailEnds(puzzle, testCoord, foundEnds)
				endCount += endsFromHere
			}
		}
	}
	return endCount
}

func part2(puzzle [][]int) int {
	total := 0

	completePathMap := make(map[Coord]int)
	for r := 0; r < len(puzzle); r++ {
		for c := 0; c < len(puzzle[0]); c++ {
			completePathMap[Coord{row: r, col: c}] = -1
		}
	}

	for r := 0; r < len(puzzle); r++ {
		for c := 0; c < len(puzzle[0]); c++ {
			if puzzle[r][c] == 0 {
				// start here
				coord := Coord{row: r, col: c}
				score := completePaths(puzzle, coord, completePathMap)
				total += score
			}
		}
	}

	return total
}

func completePaths(puzzle [][]int, coord Coord, completePathMap map[Coord]int) int {

	if completePathMap[coord] >= 0 {
		// Do we already know how many complete paths go out from this coord?
		return completePathMap[coord]
	} else {
		target := puzzle[coord.row][coord.col] + 1

		pathCount := 0
		for _, dir := range directions {
			testCoord := Coord{row: coord.row + dir.row, col: coord.col + dir.col}
			if testCoord.row >= 0 && testCoord.row < len(puzzle) && testCoord.col >= 0 && testCoord.col < len(puzzle[0]) {
				if puzzle[testCoord.row][testCoord.col] == target {
					if target == pathEndVal {
						pathCount++
					} else {
						pathsFromHere := completePaths(puzzle, testCoord, completePathMap)
						pathCount += pathsFromHere
					}
				}
			}
		}
		completePathMap[coord] = pathCount
		return pathCount
	}

}
