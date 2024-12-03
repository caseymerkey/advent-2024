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

	numsRE := regexp.MustCompile(`mul\((\d+),(\d+)\)`)

	total := 0
	for _, line := range puzzle {
		matches := numsRE.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			n1, _ := strconv.Atoi(match[1])
			n2, _ := strconv.Atoi(match[2])
			total += n1 * n2
		}

	}
	return total
}

func part2(puzzle []string) int {
	numsRE := regexp.MustCompile(`(do(?:n't)?)|(mul)\((\d+),(\d+)\)`)

	total := 0
	do := true
	for _, line := range puzzle {
		matches := numsRE.FindAllStringSubmatch(line, -1)

		for _, match := range matches {

			var op string
			if match[1] != "" {
				op = match[1]
			} else {
				op = match[2]
			}
			switch op {
			case "do":
				do = true
			case "don't":
				do = false
			case "mul":
				if do {
					n1, _ := strconv.Atoi(match[3])
					n2, _ := strconv.Atoi(match[4])
					total += n1 * n2
				}
			}
		}
	}

	return total
}
