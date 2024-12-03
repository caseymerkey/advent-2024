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
	inputFile := "input.txt"
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
		nums := make([]int, len(str))
		for i, s := range str {
			nums[i], _ = strconv.Atoi(s)
		}
		if testSafety(nums, 0) {
			safeCount++
		}
	}
	return safeCount
}

func part2(puzzle []string) int {

	safeCount := 0

	for _, line := range puzzle {

		str := strings.Fields(line)
		nums := make([]int, len(str))
		for i, s := range str {
			nums[i], _ = strconv.Atoi(s)
		}

		if testSafety(nums, 1) {
			safeCount++
		}
	}

	return safeCount
}

func testSafety(nums []int, skipTolerance int) bool {
	var previous int
	direction := 0
	safe := true

	for i, n := range nums {
		if i == 0 {
			previous = n
		} else {
			diff := previous - n
			if n != previous && diff <= 3 && diff >= -3 && direction*diff >= 0 {
				// This is okay
				direction = diff
				previous = n
			} else {
				// try removing something
				if skipTolerance == 0 {
					safe = false
					break
				} else {
					// try removing the grand-predecessor and check that
					var nums2 []int
					safe2 := false
					if i > 1 {
						nums2 = removeIndex(nums, i-2)
						safe2 = testSafety(nums2, skipTolerance-1)
					}

					// if that didn't work try removing the predecessor
					if !safe2 && i > 0 {
						nums2 = removeIndex(nums, i-1)
						safe2 = testSafety(nums2, skipTolerance-1)
					}

					if !safe2 {
						nums2 = removeIndex(nums, i)
						safe2 = testSafety(nums2, skipTolerance-1)
					}

					safe = safe2
					break
				}
			}
		}
	}

	return safe
}

func removeIndex(s []int, index int) []int {
	ret := make([]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}
