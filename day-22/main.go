package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Sequence [4]int

func (s *Sequence) push(n int) {
	for i := 0; i < 3; i++ {
		s[i] = s[i+1]
	}
	s[3] = n
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
	codes := make([]string, 0)
	for scanner.Scan() {
		codes = append(codes, scanner.Text())
	}

	startTime := time.Now()
	result := part1(codes)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	startTime = time.Now()
	result = part2(codes)
	fmt.Printf("Part 1: %d\n", result)
	executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 2 in %f seconds\n\n", executionTime)
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

func part2(seeds []string) int {
	allData := make([]map[Sequence]int, 0)
	uniqueSequences := make(map[Sequence]bool)

	for _, seed := range seeds {
		var previousPrice int
		var price int
		sequence := Sequence{}
		cache := make(map[Sequence]int)
		secret, _ := strconv.Atoi(seed)
		price = secret % 10

		for i := range 2000 {
			evaluate(&secret)
			previousPrice = price
			price = secret % 10
			delta := price - previousPrice
			sequence.push(delta)
			if i >= 3 {
				_, found := cache[sequence]
				if !found {
					cache[sequence] = price
				}
				uniqueSequences[sequence] = true
			}
		}
		allData = append(allData, cache)
	}

	highest := 0
	var bestSeq Sequence
	for seq := range uniqueSequences {
		total := 0
		for _, cache := range allData {
			total += cache[seq]
		}
		if total > highest {
			highest = total
			bestSeq = seq
		}
	}
	// fmt.Printf("%v ==> %d\n", bestSeq, highest)
	return highest
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
