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

type Coord struct {
	row int
	col int
}

func (st State) advance() State {
	d := st.facing
	loc := Coord{row: st.loc.row + d.row, col: st.loc.col + d.col}
	return State{loc: loc, facing: d}
}

type State struct {
	loc    Coord
	facing Direction
}

type Direction struct {
	row int
	col int
}

func (st State) turnRight() State {
	var facing Direction
	switch st.facing {
	case EAST:
		facing = SOUTH
	case SOUTH:
		facing = WEST
	case WEST:
		facing = NORTH
	default:
		facing = EAST
	}
	return State{loc: st.loc, facing: facing}
}
func (st State) turnLeft() State {
	var facing Direction
	switch st.facing {
	case EAST:
		facing = NORTH
	case SOUTH:
		facing = EAST
	case WEST:
		facing = SOUTH
	default:
		facing = WEST
	}
	return State{loc: st.loc, facing: facing}
}

var EAST = Direction{row: 0, col: 1}
var SOUTH = Direction{row: 1, col: 0}
var WEST = Direction{row: 0, col: -1}
var NORTH = Direction{row: -1, col: 0}

var Directions = []Direction{EAST, SOUTH, WEST, NORTH}

type Item struct {
	value    State // The value of the item
	priority int   // The priority of the item
	index    int   // The index of the item in the heap
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// For a min-heap, use pq[i].priority < pq[j].priority
	// For a max-heap, use pq[i].priority > pq[j].priority
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
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type Parents struct {
	cost    int
	parents []State
}

func (parents *Parents) evalPotentialParent(p State, cost int) {
	if parents.cost > cost {
		// new one is cheapter
		parents.parents = []State{p}
		parents.cost = cost
	} else if parents.cost == cost {
		// same cost. it's an alternative
		parents.parents = append(parents.parents, p)
	}
}

type Solver struct {
	maze         [][]string
	pq           PriorityQueue
	visitedCosts map[State]int
	parents      map[State]*Parents
	start        State
	target       Coord
	endStates    []State
	endCost      int
}

func (s Solver) isWall(c Coord) bool {
	return s.maze[c.row][c.col] == "#"
}
func (s *Solver) evaluate(current, previous State, cost int) {

	parents := s.parents[current]
	if parents == nil {
		parents = &Parents{cost: cost}
		s.parents[current] = parents
	}

	// check to see if this is a better or equal parent
	parents.evalPotentialParent(previous, cost)

	curCost, found := s.visitedCosts[current]
	if !found || cost < curCost {
		s.visitedCosts[current] = cost
		heap.Push(&s.pq, &Item{value: current, priority: cost})
	}

}

func newSolver(puzzle [][]string) Solver {

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	s := Solver{
		maze:         puzzle,
		pq:           pq,
		visitedCosts: make(map[State]int),
		parents:      make(map[State]*Parents),
		endStates:    make([]State, 0),
		endCost:      math.MaxInt,
	}

	for r, row := range s.maze {
		for c, cell := range row {
			switch cell {
			case "S":
				s.start = State{loc: Coord{row: r, col: c}, facing: EAST}
			case "E":
				s.target = Coord{row: r, col: c}
			}
		}
	}

	item := Item{value: s.start, priority: 0}
	heap.Push(&s.pq, &item)
	s.visitedCosts[s.start] = 0
	s.parents[s.start] = &Parents{cost: 0}

	return s
}

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

	r1, r2 := solve(puzzle)

	fmt.Printf("Part 1: %d\n", r1)
	fmt.Printf("Part 2: %d\n", r2)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed in %f seconds\n\n", executionTime)
}

func solve(puzzle [][]string) (int, int) {

	s := newSolver(puzzle)
	for s.pq.Len() > 0 {
		item := heap.Pop(&s.pq).(*Item)
		cheapest := item.priority
		v := item.value
		if v.loc == s.target {
			if cheapest < s.endCost {
				s.endCost = cheapest
				s.endStates = append(s.endStates, v)
			}
		} else {
			straight := v.advance()
			if !s.isWall(straight.loc) {
				s.evaluate(straight, v, cheapest+1)
			}
			left := v.turnLeft()
			s.evaluate(left, v, cheapest+1000)
			right := v.turnRight()
			s.evaluate(right, v, cheapest+1000)
		}
	}

	coordsInPath := make(map[Coord]bool)
	statesQueue := make([]State, 0)
	statesQueue = append(statesQueue, s.endStates...)
	empty := State{}

	for len(statesQueue) > 0 {
		state := statesQueue[0]
		statesQueue = statesQueue[1:]
		if state != empty {
			coordsInPath[state.loc] = true
			statesQueue = append(statesQueue, s.parents[state].parents...)
		}
	}

	return s.endCost, len(coordsInPath)
}
