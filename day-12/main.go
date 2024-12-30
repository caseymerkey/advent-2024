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
		fields := strings.Split(line, "")
		puzzle = append(puzzle, fields)
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

func part1(puzzle [][]string) int {
	total := 0

	allRegions, allPerimeters := findAllRegions(puzzle)
	for i, region := range allRegions {
		area := len(region)
		perimeter := allPerimeters[i]
		total += (area * perimeter)
	}

	return total
}

type Coord struct {
	x int
	y int
}

var directions = []Coord{
	{-1, 0},
	{0, -1},
	{1, 0},
	{0, 1},
}

func findAllRegions(puzzle [][]string) ([][]Coord, []int) {
	allRegions := make([][]Coord, 0)
	allPerimeters := make([]int, 0)
	visited := make(map[Coord]bool)

	var currentRegion []Coord

	for y, row := range puzzle {
		for x, letter := range row {
			start := Coord{x: x, y: y}
			perimeter := 0
			if !visited[start] {
				currentRegion = []Coord{start}
				visited[start] = true
				expandRegion(start, letter, &currentRegion, &perimeter, puzzle, visited)
				allRegions = append(allRegions, currentRegion)
				allPerimeters = append(allPerimeters, perimeter)
			}
		}
	}

	return allRegions, allPerimeters
}

func expandRegion(start Coord, letter string, currentRegion *[]Coord, perimeter *int, puzzle [][]string, visited map[Coord]bool) {

	for _, d := range directions {
		adj := Coord{x: start.x + d.x, y: start.y + d.y}
		if adj.x >= 0 && adj.x < len(puzzle[0]) && adj.y >= 0 && adj.y < len(puzzle) {
			if puzzle[adj.y][adj.x] == letter {
				if !visited[adj] {
					*currentRegion = append(*currentRegion, adj)
					visited[adj] = true
					expandRegion(adj, letter, currentRegion, perimeter, puzzle, visited)
				}
			} else {
				*perimeter++
			}
		} else {
			*perimeter++
		}
	}
}

func part2(puzzle [][]string) int {
	total := 0
	return total
}
