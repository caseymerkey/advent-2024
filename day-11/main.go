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
	var puzzle []string
	for scanner.Scan() {
		line := scanner.Text()
		puzzle = strings.Fields(line)
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

	workingCopy := make([]string, 0)
	workingCopy = append(workingCopy, puzzle...)

	for n := 0; n < 25; n++ {
		for i := len(workingCopy) - 1; i >= 0; i-- {
			if workingCopy[i] == "0" {
				workingCopy[i] = "1"
			} else if len(workingCopy[i])%2 == 0 {
				n1, n2 := split(workingCopy[i])
				workingCopy = slices.Insert(workingCopy, i, strconv.Itoa(n1))
				workingCopy[i+1] = strconv.Itoa(n2)

			} else {
				n, _ := strconv.Atoi(workingCopy[i])
				workingCopy[i] = strconv.Itoa(2024 * n)
			}
		}
	}
	return len(workingCopy)
}

func split(number string) (int, int) {
	n1, _ := strconv.Atoi(number[:len(number)/2])
	n2, _ := strconv.Atoi(number[len(number)/2:])
	return n1, n2
}

func part2(puzzle []string) int {
	return 0
}
