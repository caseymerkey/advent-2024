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
	codes := make([]string, 0)
	for scanner.Scan() {
		codes = append(codes, scanner.Text())
	}

	startTime := time.Now()
	result := evaluate(codes, 2)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	startTime = time.Now()
	result = evaluate(codes, 25)
	fmt.Printf("Part 2: %d\n", result)
	executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 2 in %f seconds\n", executionTime)
}

type Coord struct {
	x int
	y int
}

var NumberPad = map[string]Coord{
	"*": {x: 0, y: 0},
	"A": {x: 2, y: 0},
	"0": {x: 1, y: 0},
	"1": {x: 0, y: 1},
	"2": {x: 1, y: 1},
	"3": {x: 2, y: 1},
	"4": {x: 0, y: 2},
	"5": {x: 1, y: 2},
	"6": {x: 2, y: 2},
	"7": {x: 0, y: 3},
	"8": {x: 1, y: 3},
	"9": {x: 2, y: 3},
}

var DirectionPad = map[string]Coord{
	"<": {x: 0, y: 0},
	"v": {x: 1, y: 0},
	">": {x: 2, y: 0},
	"*": {x: 0, y: 1},
	"^": {x: 1, y: 1},
	"A": {x: 2, y: 1},
}

var DirectionButtonMap = map[Coord]string{
	{x: 0, y: 1}:  "^",
	{x: 1, y: 0}:  ">",
	{x: 0, y: -1}: "v",
	{x: -1, y: 0}: "<",
}

func evaluate(codes []string, iterations int) int {
	total := 0

	for _, code := range codes {

		cache := make(map[string]int)
		start := NumberPad["A"]
		for numberKey := range strings.SplitSeq(code, "") {
			sequence := legalPathFromAtoB(start, NumberPad[numberKey], NumberPad["*"])
			cache[sequence] = cache[sequence] + 1
			start = NumberPad[numberKey]
		}

		for range iterations {
			newCache := make(map[string]int)
			for topSeq := range cache {
				factor := cache[topSeq]
				start := DirectionPad["A"]
				for directionKey := range strings.SplitSeq(topSeq, "") {
					sequence := legalPathFromAtoB(start, DirectionPad[directionKey], DirectionPad["*"])
					newCache[sequence] = newCache[sequence] + factor
					start = DirectionPad[directionKey]
				}
			}
			cache = newCache
		}
		length := 0
		for sequence := range cache {
			length += len(sequence) * cache[sequence]
		}
		val, _ := strconv.Atoi(code[:len(code)-1])
		codeTotal := length * val
		fmt.Printf("%s ==> %d * %d = %d\n", code, val, length, codeTotal)
		total += codeTotal
	}

	return total
}

func legalPathFromAtoB(a, b, illegal Coord) string {
	var path strings.Builder

	dx := b.x - a.x
	dy := b.y - a.y
	cx := ">"
	if dx < 0 {
		cx = "<"
	}
	cy := "^"
	if dy < 0 {
		cy = "v"
	}

	// Using details posted here to decide on which order (left/right/up/down) to do first:
	// https://www.reddit.com/r/adventofcode/comments/1hjgyps/2024_day_21_part_2_i_got_greedyish
	//
	// if the move is either up-left or down-left
	// and we won't have to move through the illegal spot
	// or we can't do the up/down first because we'd be moving into the illegal space
	if (dx <= 0 && !(a.y == illegal.y && b.x == illegal.x)) || (dx > 0 && a.x == illegal.x) {
		// do left/right first, then up/down
		path.WriteString(strings.Repeat(cx, abs(dx)))
		path.WriteString(strings.Repeat(cy, abs(dy)))
	} else {
		path.WriteString(strings.Repeat(cy, abs(dy)))
		path.WriteString(strings.Repeat(cx, abs(dx)))
	}

	path.WriteString("A")
	return path.String()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
