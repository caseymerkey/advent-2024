package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"time"
)

func main() {
	inputFile := "input.txt"
	iterations := 100

	if len(os.Args) > 1 && len(os.Args[1]) > 0 {
		inputFile = os.Args[1]
	}
	if len(os.Args) > 2 {
		iterations, _ = strconv.Atoi(os.Args[2])
	}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	puzzle := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		puzzle = append(puzzle, line)
	}

	var startTime = time.Now()
	result := part1(puzzle, iterations)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	startTime = time.Now()
	result = part2(puzzle)
	fmt.Printf("Part 2: %d\n", result)
	executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 2 in %f seconds\n", executionTime)

	fmt.Println()
	part1(puzzle, 1122)
}

func part1(puzzle []string, iterations int) int {
	quadrantTotals := make([]int, 4)

	xSize := 101
	ySize := 103

	midX := (xSize - 1) / 2
	midY := (ySize - 1) / 2

	re := regexp.MustCompile(`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`)
	guards := make([]Guard, 0)
	for _, str := range puzzle {
		vals := re.FindStringSubmatch(str)
		x, _ := strconv.Atoi(vals[1])
		y, _ := strconv.Atoi(vals[2])
		dx, _ := strconv.Atoi(vals[3])
		dy, _ := strconv.Atoi(vals[4])

		endX := teleport(x+(iterations*dx), xSize)
		endY := teleport(y+(iterations*dy), ySize)

		if endX < midX && endY < midY {
			quadrantTotals[0] = quadrantTotals[0] + 1
		} else if endX < midX && endY > midY {
			quadrantTotals[1] = quadrantTotals[1] + 1
		} else if endX > midX && endY > midY {
			quadrantTotals[2] = quadrantTotals[2] + 1
		} else if endX > midX && endY < midY {
			quadrantTotals[3] = quadrantTotals[3] + 1
		}

		guards = append(guards, Guard{position: Coord{x: endX, y: endY}, direction: Coord{x: dx, y: dy}})
	}

	total := 1
	for _, n := range quadrantTotals {
		total = total * n
	}
	return total

}

func teleport(position int, boardSize int) int {

	var final int
	if position > 0 {
		final = position % boardSize
	} else {
		final = boardSize - (abs(position) % boardSize)
	}

	if final == boardSize {
		final = 0
	}
	return final
}

func abs(n int) int {
	if n < 0 {
		return -1 * n
	} else {
		return n
	}
}

type Coord struct {
	x int
	y int
}

type Guard struct {
	position  Coord
	direction Coord
}

func part2(puzzle []string) int {
	xSize := 101
	ySize := 103

	guards := make([]Guard, 0)

	re := regexp.MustCompile(`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`)

	for _, str := range puzzle {
		vals := re.FindStringSubmatch(str)
		x, _ := strconv.Atoi(vals[1])
		y, _ := strconv.Atoi(vals[2])
		dx, _ := strconv.Atoi(vals[3])
		dy, _ := strconv.Atoi(vals[4])

		g := Guard{position: Coord{x: x, y: y}, direction: Coord{x: dx, y: dy}}
		guards = append(guards, g)
	}

	for i := 1; i <= 6703; i++ {

		allLocations := make([]Coord, 0)
		for _, g := range guards {
			x := g.position.x
			y := g.position.y
			dx := g.direction.x
			dy := g.direction.y

			loc := calculatePosition(x, y, dx, dy, xSize, ySize, i)
			allLocations = append(allLocations, loc)

		}
		sequenceSize := highestSeqential(allLocations)

		if sequenceSize > 10 {
			return i
		}

	}
	return 0
}

func highestSeqential(coords []Coord) int {

	slices.SortFunc(coords, func(a Coord, b Coord) int {
		val := a.x - b.x
		if val == 0 {
			val = a.y - b.y
		}
		return val
	})

	prevX := -1
	prevY := -1
	sequenceSize := 0
	highestSequence := 0
	for _, c := range coords {

		if c.x == prevX && c.y == prevY+1 {
			prevY = c.y
			sequenceSize++
		} else {
			prevX = c.x
			prevY = c.y
			sequenceSize = 0
		}
		if sequenceSize > highestSequence {
			highestSequence = sequenceSize
		}
	}

	slices.SortFunc(coords, func(a Coord, b Coord) int {
		val := a.y - b.y
		if val == 0 {
			val = a.x - b.x
		}
		return val
	})
	prevX = -1
	prevY = -1
	sequenceSize = 0
	for _, c := range coords {

		if c.y == prevY && c.x == prevX+1 {
			prevX = c.x
			sequenceSize++
		} else {
			prevX = c.x
			prevY = c.y
			sequenceSize = 0
		}
		if sequenceSize > highestSequence {
			highestSequence = sequenceSize
		}
	}

	return highestSequence
}

func calculatePosition(x, y, dx, dy, xSize, ySize, moves int) Coord {

	endX := teleport(x+(moves*dx), xSize)
	endY := teleport(y+(moves*dy), ySize)

	return Coord{x: endX, y: endY}
}
