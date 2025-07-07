package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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
	codes := make([]string, 0)
	for scanner.Scan() {
		codes = append(codes, scanner.Text())
	}

	startTime := time.Now()
	result := part1(codes)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)
}

func part1(codes []string) int {
	total := 0
	for _, code := range codes {
		secret, _ := strconv.Atoi(code)
		for range 2000 {
			evaluate(&secret)
		}
		//fmt.Printf("%s --> %d\n", code, secret)
		total += secret
	}
	return total
}

func evaluate(secret *int) {
	val := *secret * 64
	mix(secret, val)
	prune(secret)

	val = *secret / 32
	mix(secret, val)
	prune(secret)

	val = *secret * 2048
	mix(secret, val)
	prune(secret)
}

func mix(secret *int, val int) {
	*secret = *secret ^ val
}

func prune(secret *int) {
	*secret = *secret % 16777216
}
