package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Puzzle struct {
	nodes map[string]bool
	gates []Gate
}

type Gate struct {
	n1       string
	n2       string
	function string
	output   string
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

	puzzle := Puzzle{
		nodes: make(map[string]bool),
		gates: make([]Gate, 0),
	}
	inGates := false
	inputRE := regexp.MustCompile(`(\w\d+): ([01])`)
	gateRE := regexp.MustCompile(`([\w\d]+) ([ANDXOR]+) ([\w\d]+) -> ([\w\d]+)`)
	for scanner.Scan() {
		line := scanner.Text()
		if !inGates {
			if len(line) == 0 {
				inGates = true
			} else {
				matches := inputRE.FindStringSubmatch(line)
				if matches != nil {
					name := matches[1]
					val := (matches[2] == "1")
					puzzle.nodes[name] = val
				}
			}
		} else {
			matches := gateRE.FindStringSubmatch(line)
			if matches != nil {
				n1 := matches[1]
				n2 := matches[3]
				n3 := matches[4]
				gate := Gate{
					n1:       n1,
					n2:       n2,
					function: matches[2],
					output:   n3,
				}
				puzzle.gates = append(puzzle.gates, gate)
			}
		}
	}

	var startTime = time.Now()
	result := part1(puzzle)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)
}

func part1(puzzle Puzzle) int {

	solutionValues := make(map[string]string)

	gates := puzzle.gates
	for len(gates) > 0 {
		for i, g := range slices.Backward(gates) {
			v1, ok1 := puzzle.nodes[g.n1]
			v2, ok2 := puzzle.nodes[g.n2]
			if ok1 && ok2 {
				// all inputs satisfied. eval the gate
				var v3 bool
				switch g.function {
				case "AND":
					v3 = v1 && v2
				case "OR":
					v3 = v1 || v2
				case "XOR":
					v3 = (v1 || v2) && !(v1 && v2)
				}
				puzzle.nodes[g.output] = v3

				if g.output[0] == 'z' {
					var valStr string
					if v3 {
						valStr = "1"
					} else {
						valStr = "0"
					}
					solutionValues[g.output] = valStr
				}

				// finally, remove this from the list of gates to eval
				gates = slices.Delete(gates, i, i+1)
			}
		}

	}

	keys := slices.Collect(maps.Keys(solutionValues))
	slices.SortFunc(keys, func(a, b string) int {
		// sort reverse order
		return strings.Compare(b, a)
	})
	var builder strings.Builder
	for _, s := range keys {
		builder.WriteString(solutionValues[s])
	}
	binaryValue := builder.String()
	// fmt.Println(binaryValue)

	val, _ := strconv.ParseInt(binaryValue, 2, 64)

	return int(val)
}
