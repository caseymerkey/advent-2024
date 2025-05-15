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
	resultString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(result)), ","), "[]")
	fmt.Printf("Part 1: %s\n", resultString)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)
}

func part1(availableTowels map[string]bool, desiredPatterns []string) any {

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
