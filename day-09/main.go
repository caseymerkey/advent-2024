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
	scanner.Split(bufio.ScanRunes)
	puzzle := make([]int, 0)
	for scanner.Scan() {
		r := scanner.Text()
		n, _ := strconv.Atoi(r)
		puzzle = append(puzzle, n)
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

func part1(puzzle []int) int {
	total := 0

	lastFileIdx := ((len(puzzle) - 1) / 2) * 2
	packed := make([]int, 0)
	currentFileID := 0

	toPackBuffer := make([]int, 0)
	availableFreeSpace := 0
	i := 0

	// fill up space with the next file from the front
	// find out how much space is available after that
	// while there is space
	//   is there anything in the "to pack buffer"? if so...
	//     fill up as much of that space as we can
	//   if not...
	//     grab the next available file from the END of the list
	//     put the right number if IDs in the toPackBuffer

	for {
		for availableFreeSpace > 0 {
			if len(toPackBuffer) > 0 {
				// move the toPackBuffer
				var n int
				n, toPackBuffer = pop(toPackBuffer)
				packed = append(packed, n)
				total += n * (len(packed) - 1)
				availableFreeSpace--
			} else {
				size := puzzle[lastFileIdx]
				for k := 0; k < size; k++ {
					toPackBuffer = append(toPackBuffer, lastFileIdx/2)
				}
				lastFileIdx = lastFileIdx - 2
			}
		}
		if i <= lastFileIdx {
			currentFileID = i / 2
			if currentFileID == 5 {
				fmt.Println("check it")
			}
			for k := 0; k < puzzle[i]; k++ {
				packed = append(packed, currentFileID)
				total += currentFileID * (len(packed) - 1)
			}
			availableFreeSpace = puzzle[i+1]
			i = i + 2
		} else {
			break
		}
	}
	for _, n := range toPackBuffer {
		packed = append(packed, n)
		total += (len(packed) - 1) * n
	}

	return total
}

func part2(puzzle []int) int {
	return 0
}

func pop(a []int) (int, []int) {
	if len(a) == 0 {
		return 0, a // Return 0 (zero value) if slice is empty
	}
	return a[0], a[1:]
}
