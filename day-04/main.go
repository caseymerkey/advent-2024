package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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
	puzzle := make([]string, 0)
	for scanner.Scan() {
		puzzle = append(puzzle, scanner.Text())
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

func part1(puzzle []string) int {

	count := 0
	puzzleBytes := make([][]string, 0)
	xmasRE := regexp.MustCompile("XMAS")
	samxRE := regexp.MustCompile("SAMX")
	for _, line := range puzzle {
		// fmt.Println(line)
		matches := xmasRE.FindAllString(line, -1)
		count += len(matches)
		matches = samxRE.FindAllString(line, -1)
		count += len(matches)
		// while we're at it, break the line down into a byte aray for the next step
		puzzleBytes = append(puzzleBytes, (strings.Split(line, "")))
	}
	// fmt.Println()

	// These puzzles are square
	size := len(puzzleBytes)

	for n := 0; n < size*2; n++ {
		var sb strings.Builder
		for col := 0; col < size; col++ {
			for row := 0; row < size; row++ {
				if row+col == n {
					sb.WriteString(puzzleBytes[row][col])
				}
			}
		}
		line := sb.String()
		// fmt.Println(line)
		matches := xmasRE.FindAllString(line, -1)
		count += len(matches)
		matches = samxRE.FindAllString(line, -1)
		count += len(matches)
	}

	// fmt.Println()
	for n := 1 - size; n < size; n++ {
		var sb strings.Builder
		for col := 1 - size; col < size; col++ {
			row := col - n
			if row >= 0 && row < size && col >= 0 && col < size {
				sb.WriteString(puzzleBytes[row][col])
			}
		}
		line := sb.String()
		// fmt.Println(line)
		matches := xmasRE.FindAllString(line, -1)
		count += len(matches)
		matches = samxRE.FindAllString(line, -1)
		count += len(matches)
	}

	// fmt.Println()
	puzzleBytes = rotate(puzzleBytes, 1)
	for row := 0; row < size; row++ {
		var sb strings.Builder
		for col := 0; col < size; col++ {
			sb.WriteString(puzzleBytes[row][col])
		}
		line := sb.String()
		// fmt.Println(line)
		matches := xmasRE.FindAllString(line, -1)
		count += len(matches)
		matches = samxRE.FindAllString(line, -1)
		count += len(matches)
	}

	return count
}

func part2(puzzle []string) int {

	return 0
}

func rotate(slice [][]string, direction int) [][]string {
	if direction == 0 {
		return slice
	}

	size := len(slice[0])

	result := make([][]string, size)
	for i := range result {
		result[i] = make([]string, size)
	}
	if direction > 0 {
		for r := 0; r < size; r++ {
			for c := 0; c < size; c++ {
				result[r][c] = slice[size-c-1][r]
			}
		}
	} else {
		for r := 0; r < size; r++ {
			for c := 0; c < size; c++ {
				result[r][c] = slice[c][size-r-1]
			}
		}
	}
	return result
}
