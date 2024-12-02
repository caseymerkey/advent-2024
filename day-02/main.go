package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	safeCount := 0

	for _, line := range puzzle {
		str := strings.Fields(line)
		var previous int
		direction := 0
		safe := true

		for i, s := range str {
			n, _ := strconv.Atoi(s)
			if i == 0 {
				previous = n
			} else {
				if n == previous {
					safe = false
					break
				}
				diff := previous - n
				if diff > 3 || diff < -3 {
					safe = false
					break
				}
				if direction == 0 {
					direction = diff
				} else if direction*diff < 0 {
					safe = false
					break
				}
				previous = n
			}

		}
		if safe {
			safeCount++
		}
	}
	return safeCount
}

func part2(puzzle []string) int {

	return 0
}
