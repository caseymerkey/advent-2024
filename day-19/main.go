package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	k := 0
	availableTowels := make(map[string]bool)
	var desiredPatterns []string
	for scanner.Scan() {
		line := scanner.Text()
		switch k {
		case 0:
			for t := range strings.SplitSeq(line, ",") {
				availableTowels[strings.Trim(t, " ")] = true
			}
		case 1:
			// skip
		default:
			desiredPatterns = append(desiredPatterns, line)
		}
		k++
	}

	var startTime = time.Now()
	result := part1(availableTowels, desiredPatterns)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	startTime = time.Now()
	result = part2(availableTowels, desiredPatterns)
	fmt.Printf("Part 2: %d\n", result)
	executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 2 in %f seconds\n\n", executionTime)
}

func part1(availableTowels map[string]bool, desiredPatterns []string) int {

	patternMemo := make(map[string]bool)

	var isPossible func(pattern string) bool
	isPossible = func(pattern string) bool {

		possible, found := patternMemo[pattern]
		if found {
			return possible
		}

		if availableTowels[pattern] {
			possible = true
		} else {
			for i := len(pattern) - 1; i >= 0; i-- {

				leftPossible := isPossible(pattern[:i])
				if leftPossible {
					possible = isPossible(pattern[i:])
				}
				if possible {
					break
				}
			}
		}

		patternMemo[pattern] = possible
		return possible

	}

	possibleCount := 0
	for _, pattern := range desiredPatterns {
		if isPossible(pattern) {
			possibleCount++
		}
	}
	return possibleCount

}

func part2(availableTowels map[string]bool, desiredPatterns []string) int {

	letters := []string{"w", "u", "b", "r", "g"}
	singleColors := make(map[string]bool)
	prefixPatternMap := make(map[string][]string)
	cache := make(map[string]int)

	for _, l1 := range letters {
		for _, l2 := range letters {
			key := l1 + l2
			prefixPatternMap[key] = make([]string, 0)
		}
	}

	for pattern := range availableTowels {
		if len(pattern) == 1 {
			singleColors[pattern] = true
		} else {
			prefix := pattern[0:2]
			prefixPatternMap[prefix] = append(prefixPatternMap[prefix], pattern)
		}
	}

	var combos func(target string) int
	combos = func(target string) int {

		if n, found := cache[target]; found {
			return n
		}

		if len(target) == 1 {
			if singleColors[target] {
				return 1
			} else {
				return 0
			}
		}

		count := 0
		if singleColors[target[0:1]] {
			count += combos(target[1:])
		}
		prefix := target[0:2]
		for _, pattern := range prefixPatternMap[prefix] {
			if strings.HasPrefix(target, pattern) {
				remain := target[len(pattern):]
				if len(remain) == 0 {
					count++
				} else {
					count += combos(remain)
				}
			}
		}
		cache[target] = count
		return count

	}

	total := 0
	for _, target := range desiredPatterns {
		total += combos(target)
	}

	return total
}
