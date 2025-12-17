package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var memPointer = int64(0)
var output = int64(0)
var relativeBase = int64(0)
var positionIndex = 0

// https://todd.ginsberg.com/post/advent-of-code/2019/day21/
var walk = "NOT A J\nNOT B T\nAND D T\nOR T J\nNOT C T\nOR T J\nAND D J\nWALK\n"
var run = "NOT C J\nAND D J\nNOT H T\nNOT T T\nOR E T\nAND T J\nNOT A T\nOR T J\nNOT B T\nNOT T T\nOR E T\nNOT T T\nOR T J\nRUN\n"

func copyArray(array []int64) []int64 {
	arrNew := make([]int64, 0)
	arrNew = append(arrNew, array...)

	// extend program's memory
	space := make([]int64, len(array))
	arrNew = append(arrNew, space...)
	return arrNew
}

func generateProgram() []int64 {
	file, err := os.Open("../../data/day21.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var input = []string{}

	for scanner.Scan() {
		input = strings.Split(scanner.Text(), ",")
	}

	var program = []int64{}

	for _, i := range input {
		j, err := strconv.Atoi(i)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		program = append(program, int64(j))
	}

	file.Close()

	return program
}

func setInput(part2 bool) int64 {
	var input = int64(0)
	if !part2 {
		input = int64(walk[positionIndex])
	} else {
		input = int64(run[positionIndex])
	}

	positionIndex++
	return input
}

func runOpcodeForParameterMode(opcode int64, opcodeIndex int64, program []int64, part2 bool) int {
	opcodeAsString := strconv.Itoa(int(opcode))

	i := 5 - len(opcodeAsString)

	padder := ""
	for x := 0; x < i; x++ {
		padder += "0"
	}

	// pad with 0
	opcodeAsString = padder + opcodeAsString

	code := string(opcodeAsString[3]) + string(opcodeAsString[4])
	p1 := string(opcodeAsString[2])
	p2 := string(opcodeAsString[1])
	p3 := string(opcodeAsString[0])

	o, _ := strconv.Atoi(code)
	opcode = int64(o)

	paramMode0, _ := strconv.Atoi(p1)
	paramMode1, _ := strconv.Atoi(p2)
	paramMode2, _ := strconv.Atoi(p3)

	var firstOperand = int64(0)
	var secondOperand = int64(0)

	var result = int64(0)

	opcodeIndex++

	// input instruction
	if opcode == 3 {

		input := setInput(part2)

		switch paramMode0 {
		case 0: // position mode
			program[program[opcodeIndex]] = input
		case 1: // immediate mode
			program[opcodeIndex] = input
		case 2: // relative mode
			program[program[opcodeIndex]+relativeBase] = input
		}
		return 2
	}

	switch paramMode0 {
	case 0: // position mode
		firstOperand = program[program[opcodeIndex]]
	case 1: // immediate mode
		firstOperand = program[opcodeIndex]
	case 2: // relative mode
		firstOperand = program[program[opcodeIndex]+relativeBase]
	default:
		fmt.Println("Unknown first param mode:", paramMode0)
	}

	// output instruction
	if opcode == 4 {
		output = firstOperand

		return 2
	}

	// relative base adjustment instruction
	if opcode == 9 {
		relativeBase += firstOperand
		return 2
	}

	opcodeIndex++

	switch paramMode1 {
	case 0: // position mode
		secondOperand = program[program[opcodeIndex]]
	case 1: // immediate mode
		secondOperand = program[opcodeIndex]
	case 2: // relative mode
		secondOperand = program[program[opcodeIndex]+relativeBase]
	default:
		fmt.Println("Unknown second param mode:", paramMode1)
	}

	opcodeIndex++

	switch opcode {
	case 1:
		result = firstOperand + secondOperand
	case 2:
		result = firstOperand * secondOperand
	case 5:
		if firstOperand != 0 {
			memPointer = secondOperand
			return 0
		}

		return 3
	case 6:
		if firstOperand == 0 {
			memPointer = secondOperand
			return 0
		}

		return 3
	case 7:
		if firstOperand < secondOperand {
			result = 1
		} else {
			result = 0
		}
	case 8:
		if firstOperand == secondOperand {
			result = 1
		} else {
			result = 0
		}
	}

	switch paramMode2 {
	case 0:
		program[program[opcodeIndex]] = result
	case 2:
		program[program[opcodeIndex]+relativeBase] = result
	}

	return 4
}

func runOpcode(opcode int64, opcodeIndex int64, program []int64) int {
	var firstOperand = program[program[opcodeIndex+1]]
	var secondOperand = program[program[opcodeIndex+2]]
	var resultIndex = program[opcodeIndex+3]

	var result = int64(0)

	switch opcode {
	case 1:
		result = firstOperand + secondOperand
	case 2:
		result = firstOperand * secondOperand
	case 5:
		if firstOperand != 0 {
			memPointer = secondOperand
			return 0
		}

		return 3
	case 6:
		if firstOperand == 0 {
			memPointer = secondOperand
			return 0
		}

		return 3
	case 7:
		if firstOperand < secondOperand {
			result = 1
		} else {
			result = 0
		}
	case 8:
		if firstOperand == secondOperand {
			result = 1
		} else {
			result = 0
		}
	}

	program[resultIndex] = result

	return 4
}

func runDiagnosticProgram(program []int64, part2 bool) {
	positionIndex = 0
	memPointer = int64(0)
	output = int64(0)
	relativeBase = int64(0)

	opcode := program[memPointer]

	opcodeJump := 0

	for {
		opcodeJump = 2

		switch opcode {
		case 1:
			fallthrough
		case 2:
			opcodeJump = runOpcode(opcode, memPointer, program)
		case 3:
			program[program[memPointer+1]] = setInput(part2)
		case 4:
			output = program[program[memPointer+1]]
		case 5:
			fallthrough
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			opcodeJump = runOpcode(opcode, memPointer, program)
		case 9:
			relativeBase += program[program[memPointer+1]]
		case 99:
			return
		default:
			opcodeJump = runOpcodeForParameterMode(opcode, memPointer, program, part2)
		}

		memPointer += int64(opcodeJump)
		opcode = program[memPointer]
	}
}

func main() {
	runDiagnosticProgram(copyArray(generateProgram()), false)
	fmt.Println("Part1:", output)

	runDiagnosticProgram(copyArray(generateProgram()), true)
	fmt.Println("Part2:", output)
}
