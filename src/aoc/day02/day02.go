package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func GenerateProgram() []int {
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

func RunProgram1(program []int) {
	var opcodeIndex = 0
	var opcode = program[opcodeIndex]

	for {
		firstOperand := 0
		secondOperand := 0
		resultIndex := 0

		switch opcode {
		case 1:
			// addition
			firstOperand = program[program[opcodeIndex+1]]
			secondOperand = program[program[opcodeIndex+2]]
			resultIndex = program[opcodeIndex+3]

			program[resultIndex] = firstOperand + secondOperand
			break
		case 2:
			// multiplication
			firstOperand = program[program[opcodeIndex+1]]
			secondOperand = program[program[opcodeIndex+2]]
			resultIndex = program[opcodeIndex+3]

			program[resultIndex] = firstOperand * secondOperand
			break
		case 99:
			// halt
			return
		}

		opcodeIndex += 4
		opcode = program[opcodeIndex]
	}
}

func RunProgram2(output int) {

	for noun := 0; noun <= 99; noun++ {

		for verb := 0; verb <= 99; verb++ {
			program := GenerateProgram()
			program[1] = noun
			program[2] = verb

			RunProgram1(program)
			if program[0] == output {
				part2 := (100 * program[1]) + program[2]
				fmt.Println("Part2", part2)
			}
		}
	}
}

func main() {

	var program1 = GenerateProgram()

	program1[1] = 12
	program1[2] = 2
	RunProgram1(program1)

	fmt.Println("Part1", program1[0])

	RunProgram2(19690720)
}
