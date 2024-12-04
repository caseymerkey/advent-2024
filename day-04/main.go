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

	// search horizontally
	count := 0
	puzzleBytes := make([][]string, 0)
	xmasRE := regexp.MustCompile("XMAS")
	samxRE := regexp.MustCompile("SAMX")
	for _, line := range puzzle {
		matches := xmasRE.FindAllString(line, -1)
		count += len(matches)
		matches = samxRE.FindAllString(line, -1)
		count += len(matches)
		// while we're at it, break the line down into a byte aray for the next step
		puzzleBytes = append(puzzleBytes, (strings.Split(line, "")))
	}

	// These puzzles are square
	size := len(puzzleBytes)

	// This builds strings representing the SE->NW diagonals
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
		matches := xmasRE.FindAllString(line, -1)
		count += len(matches)
		matches = samxRE.FindAllString(line, -1)
		count += len(matches)
	}

	// And NE->SW diagonals
	for n := 1 - size; n < size; n++ {
		var sb strings.Builder
		for col := 1 - size; col < size; col++ {
			row := col - n
			if row >= 0 && row < size && col >= 0 && col < size {
				sb.WriteString(puzzleBytes[row][col])
			}
		}
		line := sb.String()
		matches := xmasRE.FindAllString(line, -1)
		count += len(matches)
		matches = samxRE.FindAllString(line, -1)
		count += len(matches)
	}

	// Lastly, rotating the puzzle 90 degrees (which I had written for another problem)
	// makes getting the vertical lines easy
	puzzleBytes = rotate(puzzleBytes, 1)
	for row := 0; row < size; row++ {
		var sb strings.Builder
		for col := 0; col < size; col++ {
			sb.WriteString(puzzleBytes[row][col])
		}
		line := sb.String()
		matches := xmasRE.FindAllString(line, -1)
		count += len(matches)
		matches = samxRE.FindAllString(line, -1)
		count += len(matches)
	}

	return count
}

func part2(puzzle []string) int {

	puzzleChars := make([][]string, 0)
	for _, line := range puzzle {
		puzzleChars = append(puzzleChars, (strings.Split(line, "")))
	}

	count := 0
	for row := 1; row < len(puzzleChars)-1; row++ {
		for col := 1; col < len(puzzleChars[0])-1; col++ {
			if puzzleChars[row][col] == "A" {
				ne := puzzleChars[row-1][col-1]
				nw := puzzleChars[row-1][col+1]
				se := puzzleChars[row+1][col-1]
				sw := puzzleChars[row+1][col+1]

				if ((ne == "M" && sw == "S") || (ne == "S" && sw == "M")) && ((nw == "M" && se == "S") || (nw == "S" && se == "M")) {
					count++
				}
			}
		}
	}

	return count
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
