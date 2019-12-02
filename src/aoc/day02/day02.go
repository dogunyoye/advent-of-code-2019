package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func generateProgram() []int {
	file, err := os.Open("../../data/day02.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var input = []string{}

	for scanner.Scan() {
		input = strings.Split(scanner.Text(), ",")
	}

	var program = []int{}

	for _, i := range input {
		j, err := strconv.Atoi(i)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		program = append(program, j)
	}

	file.Close()

	return program
}

func runOpcode(opcode int, opcodeIndex int, program []int) {
	firstOperand := program[program[opcodeIndex+1]]
	secondOperand := program[program[opcodeIndex+2]]
	resultIndex := program[opcodeIndex+3]

	if opcode == 1 {
		program[resultIndex] = firstOperand + secondOperand
	} else {
		program[resultIndex] = firstOperand * secondOperand
	}
}

func runProgram1(program []int) {
	opcodeIndex := 0
	opcode := program[opcodeIndex]

	for {

		switch opcode {
		case 1:
			fallthrough
		case 2:
			runOpcode(opcode, opcodeIndex, program)
			break
		case 99:
			// halt
			return
		}

		opcodeIndex += 4
		opcode = program[opcodeIndex]
	}
}

func runProgram2(output int) {

	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			program := generateProgram()
			program[1] = noun
			program[2] = verb

			runProgram1(program)
			if program[0] == output {
				part2 := (100 * program[1]) + program[2]
				fmt.Println("Part2", part2)
			}
		}
	}
}

func main() {

	var program1 = generateProgram()

	program1[1] = 12
	program1[2] = 2
	runProgram1(program1)

	fmt.Println("Part1", program1[0])

	runProgram2(19690720)
}
