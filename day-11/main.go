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

type Node struct {
	value      string
	childCount int
	left       *Node
	right      *Node
}

var nodeMap = make(map[string]Node)

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
	var puzzle []string
	for scanner.Scan() {
		line := scanner.Text()
		puzzle = strings.Fields(line)
	}

	node := Node{value: "0"}
	nodeMap["0"] = node

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

	fmt.Println(result)
}

func part1(puzzle []string) int {

	total := 0
	for _, str := range puzzle {
		exploded := expand(str, 25)
		total += len(exploded)
	}

	return total
}

func part2(puzzle []string) int {

	total := 0
	for _, str := range puzzle {
		exploded := expand(str, 75)
		total += len(exploded)
	}

	return total
}

func expand(number string, iterations int) []string {

	result := make([]string, 0)

	node, _ := nodeMap[number]
	if node.childCount == 0 {
		if number == "0" {
			value := "1"
			left, found := nodeMap[value]
			if !found {
				left = Node{value: value}
				nodeMap[value] = left
			}
			node.left = &left
			node.childCount = 1

		} else if len(number)%2 == 0 {
			n1, n2 := split(number)
			var left, right Node
			var found bool
			value := strconv.Itoa(n1)
			left, found = nodeMap[value]
			if !found {
				left = Node{value: value}
				nodeMap[value] = left
			}
			value = strconv.Itoa(n2)
			right, found = nodeMap[value]
			if !found {
				right = Node{value: value}
				nodeMap[value] = right
			}
			node.childCount = 2
			node.left = &left
			node.right = &right
		} else {
			n, _ := strconv.Atoi(number)
			value := strconv.Itoa(2024 * n)
			left, found := nodeMap[value]
			if !found {
				left = Node{value: value}
				nodeMap[value] = left
			}
			node.left = &left
			node.childCount = 1
		}
	}

	if iterations > 1 {
		result = append(result, expand(node.left.value, iterations-1)...)
		if node.childCount > 1 {
			result = append(result, expand(node.right.value, iterations-1)...)
		}
	} else {
		result = append(result, node.left.value)
		if node.childCount > 1 {
			result = append(result, node.right.value)
		}
	}

	return result
}

func split(number string) (int, int) {
	n1, _ := strconv.Atoi(number[:len(number)/2])
	n2, _ := strconv.Atoi(number[len(number)/2:])
	return n1, n2
}
