package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
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
		line := scanner.Text()
		puzzle = append(puzzle, line)
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
	total := 0
	aButtonCost := 3
	bButtonCost := 1
	re := regexp.MustCompile(`X[\+=]([0-9]+), Y[\+=]([0-9]+)`)
	lineCounter := 0
	var xaSpaces, yaSpaces, xbSpaces, ybSpaces, xTarget, yTarget int
	for _, line := range puzzle {
		matches := re.FindStringSubmatch(line)
		switch lineCounter {
		case 0:
			xaSpaces, _ = strconv.Atoi(matches[1])
			yaSpaces, _ = strconv.Atoi(matches[2])
		case 1:
			xbSpaces, _ = strconv.Atoi(matches[1])
			ybSpaces, _ = strconv.Atoi(matches[2])
		case 2:
			xTarget, _ = strconv.Atoi(matches[1])
			yTarget, _ = strconv.Atoi(matches[2])

			aPushCount, bPushCount := findPushCounts(xaSpaces, yaSpaces, xbSpaces, ybSpaces, xTarget, yTarget)
			total += (aButtonCost * aPushCount) + (bButtonCost * bPushCount)
		}
		lineCounter++
		if lineCounter == 4 {
			lineCounter = 0
		}
	}

	return total
}

func findPushCounts(xaSpaces, yaSpaces, xbSpaces, ybSpaces, xTarget, yTarget int) (int, int) {

	var aPushCount, bPushCount int

	numerator := (xaSpaces * yTarget) - (yaSpaces * xTarget)
	denominator := (xaSpaces * ybSpaces) - (yaSpaces * xbSpaces)

	if numerator%denominator == 0 {
		bPushCount = numerator / denominator
		numerator = (xTarget - (xbSpaces * bPushCount))

		if (xTarget-(xbSpaces*bPushCount))%xaSpaces == 0 {
			aPushCount = (xTarget - (xbSpaces * bPushCount)) / xaSpaces
		}
	}
	return aPushCount, bPushCount
}

func part2(puzzle []string) int {
	total := 0
	return total
}
