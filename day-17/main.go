package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Computer struct {
	a, b, c int64
	program []int
	pointer int
	output  []int
}

func (c Computer) Run(A, B, C int64, program []int) []int {

	c.a = A
	c.b = B
	c.c = C
	c.program = program
	c.output = make([]int, 0)
	c.pointer = 0

	for c.pointer = 0; c.pointer < len(c.program); {
		c.doInstruction()
	}
	return c.output

}

func (c *Computer) doInstruction() {
	opcode := c.program[c.pointer]
	operand := c.program[c.pointer+1]

	c.pointer += 2
	switch opcode {
	case 0:
		numerator := c.a
		denominator := math.Pow(2, float64(c.comboOperandValue(operand)))
		c.a = numerator / int64(denominator)

	case 1:
		c.b = c.b ^ int64(operand)

	case 2:
		c.b = c.comboOperandValue(operand) % 8

	case 3:
		if c.a != 0 {
			c.pointer = operand
		}

	case 4:
		c.b = c.b ^ c.c

	case 5:
		c.output = append(c.output, int(c.comboOperandValue(operand)%8))

	case 6:
		numerator := c.a
		denominator := math.Pow(2, float64(c.comboOperandValue(operand)))
		c.b = numerator / int64(denominator)

	case 7:
		numerator := c.a
		denominator := math.Pow(2, float64(c.comboOperandValue(operand)))
		c.c = numerator / int64(denominator)
	}
}

func (c *Computer) comboOperandValue(op int) int64 {

	switch op {
	case 4:
		return c.a
	case 5:
		return c.b
	case 6:
		return c.c
	case 7:
		log.Fatal("Illegal Combo Operand")
	}
	return int64(op)
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
	resultString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(result)), ","), "[]")
	fmt.Printf("Part 1: %s\n", resultString)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)

	startTime = time.Now()
	result2 := part2(program)
	fmt.Printf("Part 2: %d\n", result2)
	executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 2 in %f seconds\n", executionTime)
}

func part1(registers map[string]int64, program []int) []int {

	computer := Computer{}
	output := computer.Run(registers["A"], registers["B"], registers["C"], program)
	return output
}

func part2(program []int) int64 {

	computer := Computer{}

	var dfs func(A int64, idx int) int64
	dfs = func(A int64, idx int) int64 {

		if idx < 0 {
			return A
		}

		for n := 0; n < 8; n++ {
			A2 := (A << 3) | int64(n)
			output := computer.Run(A2, 0, 0, program)
			if slices.Equal(output, program[idx:]) {
				result := dfs(A2, idx-1)
				if result > 0 {
					return result
				}
			}
		}

		return -1
	}

	return dfs(0, len(program)-1)
}
