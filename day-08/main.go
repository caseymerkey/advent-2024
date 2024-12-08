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

	for ch, arr := range antennaMap {
		fmt.Printf("%s: %v\n", ch, arr)
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
						nodes := findAntinodes(arr[a], arr[b])
						for _, n := range nodes {
							if n.row >= 0 && n.row < len(puzzle) && n.col >= 0 && n.col < len(puzzle[0]) && !slices.Contains(antinodes, n) {
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

func findAntinodes(p1 Coord, p2 Coord) []Coord {
	dx := p2.col - p1.col
	dy := p2.row - p1.row

	nodes := make([]Coord, 2)
	nodes[0] = Coord{row: p2.row + dy, col: p2.col + dx}
	nodes[1] = Coord{row: p1.row - dy, col: p1.col - dx}
	return nodes
}

func part2(puzzle [][]string) int {
	return 0
}
