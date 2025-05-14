package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultGridMax = 70
const defaultCycleCount = 1024

type Coord struct {
	x int
	y int
}

type Item struct {
	Value    any
	Priority int
}

type PriorityQueue struct {
	heap     []*Item
	indexMap map[any]int
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		heap:     []*Item{},
		indexMap: make(map[any]int),
	}
}

func (pq *PriorityQueue) Length() int {
	return len(pq.heap)
}

func (pq *PriorityQueue) Push(item *Item) {
	if _, exists := pq.indexMap[item.Value]; exists {
		return
	}
	pq.heap = append(pq.heap, item)
	pq.indexMap[item.Value] = len(pq.heap) - 1
	pq.heapifyUp(len(pq.heap) - 1)
}

func (pq *PriorityQueue) Pop() *Item {
	if pq.IsEmpty() {
		return nil
	}
	top := pq.heap[0]
	delete(pq.indexMap, top.Value)
	if len(pq.heap) > 1 {
		pq.heap[0] = pq.heap[len(pq.heap)-1]
		pq.indexMap[pq.heap[0].Value] = 0
		pq.heap = pq.heap[:len(pq.heap)-1]
		pq.heapifyDown(0)
	} else {
		pq.heap = pq.heap[:0]
	}
	return top
}

func (pq *PriorityQueue) IsEmpty() bool {
	return len(pq.heap) == 0
}

func (pq *PriorityQueue) heapifyUp(index int) {
	for index > 0 {
		parentIndex := (index - 1) / 2
		if pq.heap[index].Priority <= pq.heap[parentIndex].Priority {
			pq.swap(index, parentIndex)
			index = parentIndex
		} else {
			break
		}
	}
}

func (pq *PriorityQueue) heapifyDown(index int) {
	for {
		leftChildIndex := 2*index + 1
		rightChildIndex := 2*index + 2
		smallestIndex := index

		if leftChildIndex < len(pq.heap) && pq.heap[leftChildIndex].Priority < pq.heap[smallestIndex].Priority {
			smallestIndex = leftChildIndex
		}

		if rightChildIndex < len(pq.heap) && pq.heap[rightChildIndex].Priority < pq.heap[smallestIndex].Priority {
			smallestIndex = rightChildIndex
		}

		if smallestIndex != index {
			pq.swap(index, smallestIndex)
			index = smallestIndex
		} else {
			break
		}
	}
}

func (pq *PriorityQueue) swap(i, j int) {
	pq.heap[i], pq.heap[j] = pq.heap[j], pq.heap[i]
	pq.indexMap[pq.heap[i].Value] = i
	pq.indexMap[pq.heap[j].Value] = j
}

var Directions = []Coord{
	{x: 0, y: 1},
	{x: 1, y: 0},
	{x: 0, y: -1},
	{x: -1, y: 0},
}

func main() {

	inputFile := "sample.txt"
	if len(os.Args) > 1 && len(os.Args[1]) > 0 {
		inputFile = os.Args[1]
	}

	gridSize := defaultGridMax + 1
	cycleCount := defaultCycleCount
	if len(os.Args) > 3 {
		gridSize, _ = strconv.Atoi(os.Args[2])
		gridSize++
		cycleCount, _ = strconv.Atoi(os.Args[3])
	}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	coordList := make([]Coord, 0)
	for scanner.Scan() {
		line := scanner.Text()
		coords := strings.Split(line, ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		coordList = append(coordList, Coord{x: x, y: y})
	}

	var startTime = time.Now()
	result := part1(coordList, gridSize, cycleCount)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

}

func part1(coordList []Coord, gridSize, cycleCount int) int {

	blockedMap := make(map[Coord]bool)
	for i := 0; i < cycleCount; i++ {
		blockedMap[coordList[i]] = true
	}

	return dijkstra(blockedMap, Coord{0, 0}, Coord{gridSize - 1, gridSize - 1})
}

func dijkstra(blocked map[Coord]bool, start, target Coord) int {

	result := 0
	pq := NewPriorityQueue()
	distances := make(map[Coord]int)
	visited := make(map[Coord]bool)

	item := &Item{
		Value:    start,
		Priority: 0,
	}
	pq.Push(item)
	distances[start] = 0

	k := 0

	for pq.Length() > 0 {
		k++
		if k%100000 == 0 {
			fmt.Printf("%d - pq is %d deep\n", k, pq.Length())
		}
		pqItem := pq.Pop()
		currentLoc := pqItem.Value.(Coord)
		if currentLoc == target {
			result = distances[currentLoc]
		} else if distances[currentLoc] == math.MaxInt {
			// remaining coords are unreachable
			log.Println("Remaining locations are unreachable")
		} else {
			for _, dir := range Directions {
				c := Coord{x: currentLoc.x + dir.x, y: currentLoc.y + dir.y}
				if c.x >= 0 && c.x <= target.x && c.y >= 0 && c.y <= target.y && !blocked[c] {
					if !visited[c] {
						distance := 1 + distances[currentLoc]
						oldDistance, found := distances[c]
						if !found {
							oldDistance = math.MaxInt
						}
						if distance < oldDistance {
							distances[c] = distance
						}
						item := &Item{
							Value:    c,
							Priority: distance,
						}
						pq.Push(item)
					}
				}

			}
		}
		visited[currentLoc] = true
	}

	return result
}
