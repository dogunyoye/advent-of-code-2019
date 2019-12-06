package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var memPointer = 0

var output = 0

func copyArray(array []int) []int {
	arrNew := make([]int, 0)
	arrNew = append(arrNew, array...)
	return arrNew
}

func generateProgram() []int {
	file, err := os.Open("../../data/day05.txt")

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

func runOpcodeForParameterMode(opcode int, opcodeIndex int, program []int) int {
	opcodeAsString := strconv.Itoa(opcode)

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

	opcode, _ = strconv.Atoi(code)
	paramMode0, _ := strconv.Atoi(p1)
	paramMode1, _ := strconv.Atoi(p2)
	paramMode2, _ := strconv.Atoi(p3)

	firstOperand := 0
	secondOperand := 0

	result := 0

	opcodeIndex++

	if paramMode0 == 0 {
		firstOperand = program[program[opcodeIndex]]
	} else {
		firstOperand = program[opcodeIndex]
	}

	if opcode == 4 {
		if paramMode0 == 0 {
			output = program[firstOperand]
		} else {
			output = firstOperand
		}
		return 2
	}

	opcodeIndex++

	if paramMode1 == 0 {
		secondOperand = program[program[opcodeIndex]]
	} else {
		secondOperand = program[opcodeIndex]
	}

	opcodeIndex++

	switch opcode {
	case 1:
		result = firstOperand + secondOperand
		break
	case 2:
		result = firstOperand * secondOperand
		break
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

	if paramMode2 == 0 {
		program[program[opcodeIndex]] = result
	}

	return 4
}

func runOpcode(opcode int, opcodeIndex int, program []int) int {
	firstOperand := program[program[opcodeIndex+1]]
	secondOperand := program[program[opcodeIndex+2]]
	resultIndex := program[opcodeIndex+3]

	result := 0

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

func runDiagnosticProgram(program []int, input int) {
	opcode := program[memPointer]

	opcodeJump := 0

	for {
		opcodeJump = 2

		switch opcode {
		case 1:
			fallthrough
		case 2:
			opcodeJump = runOpcode(opcode, memPointer, program)
			break
		case 3:
			program[program[memPointer+1]] = input
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
			break
		case 99:
			return
		default:
			opcodeJump = runOpcodeForParameterMode(opcode, memPointer, program)
		}

		memPointer += opcodeJump
		opcode = program[memPointer]
	}
}

func main() {

	var program1 = copyArray(generateProgram())
	//var program2 = copyArray(program1)

	runDiagnosticProgram(program1, 1)
	fmt.Println("Part1:", output)
	memPointer = 0

	// runDiagnosticProgram(program2, 5)
	// fmt.Println("Part2:", output)
}
