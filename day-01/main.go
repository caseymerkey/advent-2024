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
	puzzle := make([]string, 0)
	for scanner.Scan() {
		puzzle = append(puzzle, scanner.Text())
	}

	var startTime = time.Now()
	p1 := part1(puzzle)
	fmt.Printf("Part 1: %d\n", p1)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n", executionTime)
}

func part1(puzzle []string) int {

	leftSide := make([]int, 0)
	rightSide := make([]int, 0)
	for _, line := range puzzle {
		s := strings.Fields(line)
		val, _ := strconv.Atoi(s[0])
		leftSide = append(leftSide, val)
		val, _ = strconv.Atoi(s[1])
		rightSide = append(rightSide, val)
	}
	slices.Sort(leftSide)
	slices.Sort(rightSide)

	total := 0
	for i, _ := range leftSide {
		left := leftSide[i]
		right := rightSide[i]
		diff := 0
		if left > right {
			diff = left - right
		} else {
			diff = right - left
		}
		total += diff
	}

	return total
}
