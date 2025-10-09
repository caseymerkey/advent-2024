package main

import (
	"bufio"
	"errors"
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

type Gate struct {
	n1       string
	n2       string
	function string
	n3       string
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

	nodes := make(map[string]bool)
	gates := make([]Gate, 0)

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
					nodes[name] = val
				}
			}
		} else {
			matches := gateRE.FindStringSubmatch(line)
			if matches != nil {
				gate := Gate{
					n1:       matches[1],
					function: matches[2],
					n2:       matches[3],
					n3:       matches[4],
				}
				gates = append(gates, gate)
			}
		}
	}

	var startTime = time.Now()
	r1 := part1(copyNodes(nodes), gates)
	fmt.Printf("Part 1: %d\n", r1)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	// startTime = time.Now()
	// r2 := part2(nodes, gates)
	// fmt.Printf("Part 2: %s\n", r2)
	// executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	// fmt.Printf("Completed Part 2 in %f seconds\n\n", executionTime)
}

func part1(nodes map[string]bool, gates []Gate) int {
	// eval all gates
	evalAllGates(nodes, gates)
	// return z values
	return readNumber(nodes, 'z')
}

func evalAllGates(nodes map[string]bool, gates []Gate) {
	for len(gates) > 0 {
		for i, g := range slices.Backward(gates) {
			v3, err := evalGate(nodes, g)
			if err == nil {
				nodes[g.n3] = v3
				gates = slices.Delete(gates, i, i+1)
			}
		}
	}
}

func readNumber(nodes map[string]bool, prefix byte) int {
	solutionValues := make(map[string]string)

	for k, v := range nodes {
		if k[0] == prefix {
			solutionValues[k] = boolToBinaryString(v)
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
	binaryString := builder.String()
	val, _ := strconv.ParseInt(binaryString, 2, 64)

	return int(val)
}

func boolToBinaryString(b bool) string {
	valStr := "0"
	if b {
		valStr = "1"
	}
	return valStr
}

func evalGate(nodes map[string]bool, g Gate) (bool, error) {
	v1, ok1 := nodes[g.n1]
	v2, ok2 := nodes[g.n2]
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
		return v3, nil
	} else {
		return false, errors.New("Inputs not available")
	}
}

func copyNodes(nodes map[string]bool) map[string]bool {

	copy := make(map[string]bool)
	for k, v := range nodes {
		copy[k] = v
	}
	return copy
}
