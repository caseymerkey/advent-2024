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
	var puzzle []string
	for scanner.Scan() {
		line := scanner.Text()
		puzzle = strings.Fields(line)
	}

	var startTime = time.Now()
	result := solve(puzzle, 25)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	startTime = time.Now()
	result = solve(puzzle, 75)
	fmt.Printf("Part 2: %d\n", result)
	executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 2 in %f seconds\n", executionTime)

	fmt.Println(result)
}

func split(number string) (int, int) {
	n1, _ := strconv.Atoi(number[:len(number)/2])
	n2, _ := strconv.Atoi(number[len(number)/2:])
	return n1, n2
}

func solve(puzzle []string, iterations int) int {

	expandMemo := make(map[string][]string)

	puzzleExpansion := make(map[string]int)
	for _, stone := range puzzle {
		puzzleExpansion[stone]++
	}

	for i := 0; i < iterations; i++ {
		newExpansion := make(map[string]int)
		for stone, count := range puzzleExpansion {
			expanded := handleSingleStone(stone, expandMemo)
			for _, st := range expanded {
				newExpansion[st] += count
			}
		}
		puzzleExpansion = newExpansion
	}

	result := 0
	for _, n := range puzzleExpansion {
		result += n
	}
	return result
}

func handleSingleStone(stone string, expandMemo map[string][]string) []string {
	result := make([]string, 0)

	r, found := expandMemo[stone]
	if found {
		result = append(result, r...)
	} else {
		if stone == "0" {
			result = append(result, "1")
		} else if len(stone)%2 == 0 {
			n1, n2 := split(stone)
			result = append(result, strconv.Itoa(n1), strconv.Itoa(n2))
		} else {
			n, _ := strconv.Atoi(stone)
			result = append(result, strconv.Itoa(n*2024))
		}
		expandMemo[stone] = result
	}
	return result
}
