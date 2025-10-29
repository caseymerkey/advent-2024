package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type KeyLock struct {
	pattern []string
	values  []int
}

func NewKeyLock(pattern []string) KeyLock {
	kl := KeyLock{pattern: pattern, values: make([]int, len(pattern[0]))}

	for _, row := range pattern {
		for c, col := range row {
			if col == '#' {
				kl.values[c]++
			}
		}
	}
	return kl
}

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
	block := make([]string, 0)

	keys := make([]KeyLock, 0)
	locks := make([]KeyLock, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			kl := NewKeyLock(block)
			if block[0] == "#####" {
				keys = append(keys, kl)
			} else {
				locks = append(locks, kl)
			}
			block = make([]string, 0)
		} else {
			block = append(block, line)
		}
	}
	kl := NewKeyLock(block)
	if block[0] == "#####" {
		keys = append(keys, kl)
	} else {
		locks = append(locks, kl)
	}

	var startTime = time.Now()
	r1 := part1(locks, keys)
	fmt.Printf("Part 1: %d\n", r1)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	// startTime = time.Now()
	// r2 := part2(locks, keys)
	// fmt.Printf("Part 2: %s\n", r2)
	// executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	// fmt.Printf("Completed Part 2 in %f seconds\n\n", executionTime)
}

func part1(locks []KeyLock, keys []KeyLock) int {

	matchCount := 0
	for _, lock := range locks {
		for _, key := range keys {
			if checkFit(lock, key) {
				matchCount++
			}
		}
	}

	return matchCount
}

func checkFit(key, lock KeyLock) bool {
	fits := true
	for i, val := range key.values {
		if val+lock.values[i] > 7 {
			fits = false
			break
		}
	}
	return fits
}
