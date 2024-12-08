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
		puzzle = append(puzzle, strings.Split(line, ""))
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

func part1(puzzle [][]string) int {

	antennaMap := make(map[string][]Coord)
	antinodes := make([]Coord, 0)

	for r, row := range puzzle {
		for c, ch := range row {
			if ch != "." {
				coords := antennaMap[ch]
				coords = append(coords, Coord{row: r, col: c})
				antennaMap[ch] = coords
			}
		}
	}

	for _, arr := range antennaMap {
		evaluated := make(map[string]bool)
		for a := 0; a < len(arr); a++ {
			for b := 0; b < len(arr); b++ {
				if a != b {
					var key string
					if a < b {
						key = strconv.Itoa(a) + "-" + strconv.Itoa(b)
					} else {
						key = strconv.Itoa(b) + "-" + strconv.Itoa(a)
					}
					if !evaluated[key] {
						evaluated[key] = true
						nodes := findAntinodes(arr[a], arr[b], len(puzzle), false)
						for _, n := range nodes {
							if !slices.Contains(antinodes, n) {
								antinodes = append(antinodes, n)
							}
						}
					}
				}
			}
		}
	}

	return len(antinodes)
}

func part2(puzzle [][]string) int {
	antennaMap := make(map[string][]Coord)
	antinodes := make([]Coord, 0)

	for r, row := range puzzle {
		for c, ch := range row {
			if ch != "." {
				coords := antennaMap[ch]
				coords = append(coords, Coord{row: r, col: c})
				antennaMap[ch] = coords
			}
		}
	}

	for _, arr := range antennaMap {
		evaluated := make(map[string]bool)
		for a := 0; a < len(arr); a++ {
			for b := 0; b < len(arr); b++ {
				if a != b {
					var key string
					if a < b {
						key = strconv.Itoa(a) + "-" + strconv.Itoa(b)
					} else {
						key = strconv.Itoa(b) + "-" + strconv.Itoa(a)
					}
					if !evaluated[key] {
						evaluated[key] = true
						nodes := findAntinodes(arr[a], arr[b], len(puzzle), true)
						for _, n := range nodes {
							if !slices.Contains(antinodes, n) {
								antinodes = append(antinodes, n)
							}
						}
					}
				}
			}
		}
	}

	for _, c := range antinodes {
		if puzzle[c.row][c.col] == "." {
			puzzle[c.row][c.col] = "#"
		}
	}
	for r := 0; r < len(puzzle); r++ {
		for c := 0; c < len(puzzle[r]); c++ {
			fmt.Print(puzzle[r][c])
		}
		fmt.Println()
	}

	return len(antinodes)
}

func findAntinodes(p1 Coord, p2 Coord, gridSize int, considerHarmonics bool) []Coord {
	dx := p2.col - p1.col
	dy := p2.row - p1.row

	nodes := make([]Coord, 0)

	factor := 1
	for factor == 1 || considerHarmonics {
		n := Coord{row: p2.row + (dy * factor), col: p2.col + (dx * factor)}
		if n.row >= 0 && n.row < gridSize && n.col >= 0 && n.col < gridSize {
			nodes = append(nodes, n)
		} else {
			break
		}
		factor++
	}

	factor = 1
	for factor == 1 || considerHarmonics {
		n := Coord{row: p1.row - (dy * factor), col: p1.col - (dx * factor)}
		if n.row >= 0 && n.row < gridSize && n.col >= 0 && n.col < gridSize {
			nodes = append(nodes, n)
		} else {
			break
		}
		factor++
	}

	if considerHarmonics {
		nodes = append(nodes, p1)
		nodes = append(nodes, p2)
	}
	return nodes
}
