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

	predecessorRules := make(map[int][]int)
	updates := make([][]int, 0)
	total := 0

	updateSection := false
	for _, line := range puzzle {
		if line == "" {
			updateSection = true
		} else {
			if updateSection {
				vals := strings.Split(line, ",")
				pages := make([]int, 0)
				for _, v := range vals {
					p, _ := strconv.Atoi(v)
					pages = append(pages, p)
				}
				updates = append(updates, pages)
			} else {
				vals := strings.Split(line, "|")
				later, _ := strconv.Atoi(vals[1])
				earlier, _ := strconv.Atoi(vals[0])
				predecessors := predecessorRules[later]
				predecessors = append(predecessors, earlier)
				predecessorRules[later] = predecessors
			}
		}
	}

	for _, update := range updates {
		valid := true
		for i, p := range update {
			laterPages := update[i+1:]
			predecessors := predecessorRules[p]

			for _, p2 := range laterPages {
				if slices.Contains(predecessors, p2) {
					valid = false
					break
				}
			}
			if !valid {
				break
			}
		}
		if valid {
			middle := update[(len(update)-1)/2]
			total += middle
		}
	}

	return total
}

func part2(puzzle []string) int {
	return 0
}
