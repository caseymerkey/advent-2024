package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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

func part1(codes []string) int {
	total := 0

	var expand func(directionButtons []string, level int) []string
	expand = func(directionButtons []string, level int) []string {
		if level == 0 {
			return directionButtons
		}

		expanded := make([]string, 0)
		start := DirectionPad["A"]

		for _, dirButton := range directionButtons {
			target := DirectionPad[dirButton]
			buttonPath := legalPathFromAtoB(start, target, DirectionPad["*"])
			expanded = append(expanded, buttonPath...)
			start = target
		}

		return expand(expanded, level-1)
	}

	for _, code := range codes {
		numPadKeys := strings.Split(code, "")
		controlKeys := make([]string, 0)
		start := NumberPad["A"]
		for _, numKey := range numPadKeys {
			target := NumberPad[numKey]
			controlKeys = append(controlKeys, legalPathFromAtoB(start, target, NumberPad["*"])...)
			start = target
		}
		controlKeys = expand(controlKeys, 2)

		l, _ := strconv.Atoi(code[:len(code)-1])
		total += len(controlKeys) * l

	}

	return total
}

func legalPathFromAtoB(a, b, illegal Coord) []string {
	path := make([]string, 0)

	dx := b.x - a.x
	dy := b.y - a.y
	cx := []string{">"}
	if dx < 0 {
		cx[0] = "<"
	}
	cy := []string{"^"}
	if dy < 0 {
		cy[0] = "v"
	}

	// Using details posted here to decide on which order (left/right/up/down) to do first:
	// https://www.reddit.com/r/adventofcode/comments/1hjgyps/2024_day_21_part_2_i_got_greedyish
	//
	// if the move is either up-left or down-left
	// and we won't have to move through the illegal spot
	// or we can't do the up/down first because we'd be moving into the illegal space
	if (dx <= 0 && !(a.y == illegal.y && b.x == illegal.x)) || (dx > 0 && a.x == illegal.x) {
		// do left/right first, then up/down
		path = append(path, slices.Repeat(cx, int(math.Abs(float64(dx))))...)
		path = append(path, slices.Repeat(cy, int(math.Abs(float64(dy))))...)
	} else {
		path = append(path, slices.Repeat(cy, int(math.Abs(float64(dy))))...)
		path = append(path, slices.Repeat(cx, int(math.Abs(float64(dx))))...)
	}

	path = append(path, "A")
	return path
}
