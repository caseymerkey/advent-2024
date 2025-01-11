package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
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
	puzzle := make([][]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		puzzle = append(puzzle, strings.Split(line, ""))
	}

	var startTime = time.Now()
	result := part1(puzzle)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	// startTime = time.Now()
	// result = part2(puzzle)
	// fmt.Printf("Part 2: %d\n", result)
	// executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	// fmt.Printf("Completed Part 2 in %f seconds\n", executionTime)
}

func part1(puzzle [][]string) int {

	var start, target Coord
	for r, row := range puzzle {
		for c, str := range row {
			if str == "S" {
				start = Coord{row: r, col: c}
			} else if str == "E" {
				target = Coord{row: r, col: c}
			}
		}
	}
	// 82464 is too high
	// return aStar(puzzle, start.row, start.col, target.row, target.col)

	return dijkstra(puzzle, start, target)
}

type Coord struct {
	row int
	col int
}
type Node struct {
	loc Coord
	//parent Coord
	facing int // 0 == left/right, 1 == up/down

}

type Item struct {
	value    Node
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Item, value Node, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

var Directions = []Coord{
	{row: 0, col: 1},
	{row: 1, col: 0},
	{row: 0, col: -1},
	{row: -1, col: 0},
}

func dijkstra(grid [][]string, start Coord, target Coord) int {

	startNode := Node{loc: start, facing: 0}

	pq := make(PriorityQueue, 0)
	distances := make(map[Node]int)
	visited := make(map[Node]bool)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    startNode,
		priority: 0,
	}
	heap.Push(&pq, item)

	distances[startNode] = 0

	for len(pq) > 0 {

		pqItem := heap.Pop(&pq).(*Item)
		currentNode := pqItem.value

		if currentNode.loc == target {
			fmt.Println(" ****  We're at the target  ****")
			fmt.Printf(" ****  Distance = %d  ****\n", distances[currentNode])
			return distances[currentNode]
		} else if distances[currentNode] == math.MaxInt {
			// these remaining nodes are unreachable
			fmt.Println(" ****  Remaining nodes are unreachable *****")
		} else {
			for _, dir := range Directions {
				c := Coord{row: currentNode.loc.row + dir.row, col: currentNode.loc.col + dir.col}

				if c.row >= 0 && c.row < len(grid) && c.col >= 0 && c.col < len(grid[0]) && grid[c.row][c.col] != "#" {

					facing := 0
					if dir.col == 0 {
						facing = 1
					}
					neighbor := Node{loc: c, facing: facing}
					v, _ := visited[neighbor]
					if !v {
						distance := 1 + distances[currentNode]
						if currentNode.facing != facing {
							distance += 1000
						}

						oldDistance, found := distances[neighbor]
						if !found {
							oldDistance = math.MaxInt
						}
						if distance < oldDistance {
							distances[neighbor] = distance
						}

						item := &Item{
							value:    neighbor,
							priority: distance,
						}
						heap.Push(&pq, item)
					}
				}
			}
		}
		visited[currentNode] = true
	}

	return 0
}
