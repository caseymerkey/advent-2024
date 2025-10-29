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

var functions = map[string]func(bool, bool) bool{
	"AND": func(v1, v2 bool) bool {
		return v1 && v2
	},
	"OR": func(v1, v2 bool) bool {
		return v1 || v2
	},
	"XOR": func(v1, v2 bool) bool {
		return (v1 || v2) && !(v1 && v2)
	},
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

	startTime = time.Now()
	r2 := part2(nodes, gates)
	fmt.Printf("Part 2: %s\n", r2)
	executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 2 in %f seconds\n\n", executionTime)
}

func part1(nodes map[string]bool, gates []Gate) int {
	// eval all gates
	evalAllGates(nodes, gates)
	// return z values
	return readNumber(nodes, 'z')
}

func part2(nodes map[string]bool, gates []Gate) string {

	//If the output of a gate is z, then the operation has to be XOR unless it is the last bit.
	//If the output of a gate is not z and the inputs are not x, y then it has to be AND / OR, but not XOR.

	//Also...
	//If you have a XOR gate with inputs x, y, there must be another XOR gate with this gate as an input. Search through all gates for an XOR-gate with this gate as an input; if it does not exist, your (original) XOR gate is faulty.
	//Similarly, if you have an AND-gate, there must be an OR-gate with this gate as an input. If that gate doesn't exist, the original AND gate is faulty.
	//(These don't apply for the gates with input x00, y00).
	// from: https://www.reddit.com/r/adventofcode/comments/1hla5ql/2024_day_24_part_2_a_guide_on_the_idea_behind_the/

	maxX := 0
	for i := range nodes {
		if i[0] == 'x' {
			n, _ := strconv.Atoi(i[1:])
			if n > maxX {
				maxX = n
			}
		}
	}
	maxZ := fmt.Sprintf("z%d", maxX+1)
	sketchyWires := make([]string, 0)

	for _, g := range gates {
		inputs := []byte{g.n1[0], g.n2[0]}
		slices.Sort(inputs)
		if g.n3[0] == 'z' {
			if g.n3 != maxZ && g.function != "XOR" {
				sketchyWires = append(sketchyWires, g.n3)
				continue
			}
		} else {
			if inputs[0] != 'x' || inputs[1] != 'y' {
				if g.function == "XOR" {
					sketchyWires = append(sketchyWires, g.n3)
					continue
				}
			}
		}
		if g.function == "XOR" && inputs[0] == 'x' && inputs[1] == 'y' && (g.n1 != "x00" && g.n1 != "y00" && g.n2 != "x00" && g.n2 != "y00") {
			foundMatch := false
			for _, g1 := range gates {
				if g1.function == "XOR" && (g1.n1 == g.n3 || g1.n2 == g.n3) {
					foundMatch = true
					break
				}
			}
			if !foundMatch {
				sketchyWires = append(sketchyWires, g.n3)
			}
			continue
		}
		if g.function == "AND" && inputs[0] == 'x' && inputs[1] == 'y' && (g.n1 != "x00" && g.n1 != "y00" && g.n2 != "x00" && g.n2 != "y00") {
			foundMatch := false
			for _, g1 := range gates {
				if g1.function == "OR" && (g1.n1 == g.n3 || g1.n2 == g.n3) {
					foundMatch = true
					break
				}
			}
			if !foundMatch {
				sketchyWires = append(sketchyWires, g.n3)
			}
			continue
		}
	}

	slices.Sort(sketchyWires)
	return fmt.Sprintf("%v\n", strings.Join(sketchyWires, ","))
}

func evalAllGates(nodes map[string]bool, gates []Gate) {
	gatesCopy := make([]Gate, len(gates))
	copy(gatesCopy, gates)
	for len(gatesCopy) > 0 {
		for i, g := range slices.Backward(gatesCopy) {
			v3, err := evalGate(nodes, g)
			if err == nil {
				nodes[g.n3] = v3
				gatesCopy = slices.Delete(gatesCopy, i, i+1)
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
	fmt.Println(binaryString)
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
		f := functions[g.function]
		v3 := f(v1, v2)
		return v3, nil
	} else {
		return false, errors.New("Inputs not available")
	}
}

func copyNodes(nodes map[string]bool) map[string]bool {

	copy := make(map[string]bool)
	maps.Copy(copy, nodes)
	return copy
}
