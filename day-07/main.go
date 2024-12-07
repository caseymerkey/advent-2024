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

type Equation struct {
	target   int
	operands []int
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
	puzzleText := make([]string, 0)
	for scanner.Scan() {
		puzzleText = append(puzzleText, scanner.Text())
	}

	puzzle := make([]Equation, 0)
	for _, line := range puzzleText {
		s := strings.Split(line, ":")
		target, _ := strconv.Atoi(s[0])
		s = strings.Fields(strings.Trim(s[1], " "))
		vals := make([]int, len(s))
		for i, txt := range s {
			vals[i], _ = strconv.Atoi(txt)
		}
		puzzle = append(puzzle, Equation{target: target, operands: vals})
	}

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
}

func part1(puzzle []Equation) int {
	total := 0
	for i, equation := range puzzle {
		calc_memo := make(map[string]int)

		if calculate("", equation.operands, equation.target, calc_memo) {
			fmt.Printf("Equation %d works\n", i)
			total += equation.target
		}
	}

	return total
}

func calculate(evaluated string, remainingOperands []int, target int, calc_memo map[string]int) bool {

	if evaluated == "" {

		a := remainingOperands[0]
		b := remainingOperands[1]
		remaining := remainingOperands[2:]

		if len(remaining) == 0 {
			return (a*b == target || a+b == target)
		}
		key := strconv.Itoa(a) + "*" + strconv.Itoa(b)
		calc_memo[key] = a * b
		if calculate(key, remaining, target, calc_memo) {
			return true
		}

		key = strconv.Itoa(a) + "+" + strconv.Itoa(b)
		calc_memo[strconv.Itoa(a)+"+"+strconv.Itoa(b)] = a + b
		if calculate(key, remaining, target, calc_memo) {
			return true
		}
		return false

	} else {
		valid := false
		// look in the memo to find value of the evaluated part
		v, ok := calc_memo[evaluated]
		if ok {

			next, remaining := pop(remainingOperands)

			for i := 0; i < 2; i++ {
				result, op := eval(v, next, i)
				key := evaluated + op + strconv.Itoa(next)
				calc_memo[key] = result
				if len(remaining) == 0 {
					if result == target {
						return true

					}
				} else {
					if result > target {
						valid = false
					} else {
						valid = calculate(key, remaining, target, calc_memo)
						if valid {
							return true
						}
					}
				}
			}

		} else {
			fmt.Println("never seen this before")
			valid = false
		}

		// evaluate previous * next operand
		// store result in memo
		// if the new result is already larger than target then return 0
		// if there are no more operands and te result == target return 1
		// otherwise if there are more operands
		// recursively call calculate with the new evaluated string & remaining operands

		// evaluate previous + next operand... etc.

		return valid
	}
}

func eval(a int, b int, operator int) (int, string) {
	if operator == 0 {
		return a * b, "*"
	} else {
		return a + b, "+"
	}
}

func pop(a []int) (int, []int) {
	if len(a) == 0 {
		return 0, a // Return 0 (zero value) if slice is empty
	}
	return a[0], a[1:]
}

func part2(puzzle []Equation) int {
	return 0
}
