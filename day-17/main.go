package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
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

	registers := make(map[string]int64)
	program := make([]int, 0)
	registerRE := regexp.MustCompile(`Register ([A-C]): ([\d]+)`)
	programRE := regexp.MustCompile(`Program: ([\d,]+)`)
	for scanner.Scan() {
		line := scanner.Text()
		if registerRE.Match([]byte(line)) {
			matches := registerRE.FindStringSubmatch(line)
			n, _ := strconv.Atoi(matches[2])
			registers[matches[1]] = int64(n)
		} else if programRE.Match([]byte(line)) {
			matches := programRE.FindStringSubmatch(line)
			for s := range strings.SplitSeq(matches[1], ",") {
				n, _ := strconv.Atoi(s)
				program = append(program, n)
			}
		}
	}

	var startTime = time.Now()
	result := part1(registers, program)
	fmt.Printf("Part 1: %s\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	// startTime = time.Now()
	// result = part2(puzzle)
	// fmt.Printf("Part 2: %d\n", result)
	// executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	// fmt.Printf("Completed Part 2 in %f seconds\n", executionTime)
}

func part1(registers map[string]int64, program []int) string {

	output := make([]string, 0)
	for pointer := 0; pointer < len(program); {
		opcode := program[pointer]
		operand := program[pointer+1]

		ptr, out := doInstruction(opcode, operand, registers)
		if ptr >= 0 {
			pointer = ptr
		} else {
			pointer += 2
		}
		if out >= 0 {
			output = append(output, strconv.Itoa(out))
		}

	}

	return strings.Join(output, ",")
}

func comboOperandValue(op int, registers map[string]int64) int64 {

	switch op {
	case 4:
		return registers["A"]
	case 5:
		return registers["B"]
	case 6:
		return registers["C"]
	case 7:
		log.Fatal("Illegal Combo Operand")
	}
	return int64(op)
}

func doInstruction(operator int, operand int, registers map[string]int64) (int, int) {

	jump := -1
	out := -1

	switch operator {
	case 0:
		numerator := registers["A"]
		denominator := math.Pow(2, float64(comboOperandValue(operand, registers)))
		registers["A"] = numerator / int64(denominator)

	case 1:
		registers["B"] = registers["B"] ^ int64(operand)

	case 2:
		registers["B"] = comboOperandValue(operand, registers) % 8

	case 3:
		if registers["A"] != 0 {
			jump = operand
		}

	case 4:
		registers["B"] = registers["B"] ^ registers["C"]

	case 5:
		out = int(comboOperandValue(operand, registers) % 8)

	case 6:
		numerator := registers["A"]
		denominator := math.Pow(2, float64(comboOperandValue(operand, registers)))
		registers["B"] = numerator / int64(denominator)

	case 7:
		numerator := registers["A"]
		denominator := math.Pow(2, float64(comboOperandValue(operand, registers)))
		registers["C"] = numerator / int64(denominator)
	}

	return jump, out
}
