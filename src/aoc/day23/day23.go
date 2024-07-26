package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type intcodeComputer struct {
	id           int
	program      []int64
	memPointer   int64
	output       int64
	relativeBase int64
	queue        []packet
}

type packet struct {
	x int64
	y int64
}

func setInput() int64 {
	var input = int64(0)
	return input
}

func (computer *intcodeComputer) runOpcodeForParameterMode(opcode int64, opcodeIndex int64) int {
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

		input := setInput()

		if paramMode0 == 0 { // position mode
			computer.program[computer.program[opcodeIndex]] = input
		} else if paramMode0 == 1 { // immediate mode
			computer.program[opcodeIndex] = input
		} else if paramMode0 == 2 { // relative mode
			computer.program[computer.program[opcodeIndex]+computer.relativeBase] = input
		}
		return 2
	}

	if paramMode0 == 0 { // position mode
		firstOperand = computer.program[computer.program[opcodeIndex]]
	} else if paramMode0 == 1 { // immediate mode
		firstOperand = computer.program[opcodeIndex]
	} else if paramMode0 == 2 { // relative mode
		firstOperand = computer.program[computer.program[opcodeIndex]+computer.relativeBase]
	} else {
		fmt.Println("Unknown first param mode:", paramMode0)
	}

	// output instruction
	if opcode == 4 {
		computer.output = firstOperand
		return 2
	}

	// relative base adjustment instruction
	if opcode == 9 {
		computer.relativeBase += firstOperand
		return 2
	}

	opcodeIndex++

	if paramMode1 == 0 { // position mode
		secondOperand = computer.program[computer.program[opcodeIndex]]
	} else if paramMode1 == 1 { // immediate mode
		secondOperand = computer.program[opcodeIndex]
	} else if paramMode1 == 2 { // relative mode
		secondOperand = computer.program[computer.program[opcodeIndex]+computer.relativeBase]
	} else {
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
			computer.memPointer = secondOperand
			return 0
		}

		return 3
	case 6:
		if firstOperand == 0 {
			computer.memPointer = secondOperand
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
		computer.program[computer.program[opcodeIndex]] = result
	} else if paramMode2 == 2 {
		computer.program[computer.program[opcodeIndex]+computer.relativeBase] = result
	}

	return 4
}

func (computer *intcodeComputer) runOpcode(opcode int64, opcodeIndex int64, program []int64) int {
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
			computer.memPointer = secondOperand
			return 0
		}

		return 3
	case 6:
		if firstOperand == 0 {
			computer.memPointer = secondOperand
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

func (computer *intcodeComputer) runDiagnosticProgram() {
	opcode := computer.program[computer.memPointer]
	opcodeJump := 0

	for {
		opcodeJump = 2

		switch opcode {
		case 1:
			fallthrough
		case 2:
			opcodeJump = computer.runOpcode(opcode, computer.memPointer, computer.program)
		case 3:
			computer.program[computer.program[computer.memPointer+1]] = setInput()
		case 4:
			computer.output = computer.program[computer.program[computer.memPointer+1]]
		case 5:
			fallthrough
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			opcodeJump = computer.runOpcode(opcode, computer.memPointer, computer.program)
		case 9:
			computer.relativeBase += computer.program[computer.program[computer.memPointer+1]]
		case 99:
			return
		default:
			opcodeJump = computer.runOpcodeForParameterMode(opcode, computer.memPointer)
		}

		computer.memPointer += int64(opcodeJump)
		opcode = computer.program[computer.memPointer]
	}
}

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

func findYValueOfFirstPacketSentToAddress255() int {
	computers := make([]*intcodeComputer, 0)
	for i := 0; i < 50; i++ {
		c := intcodeComputer{i, copyArray(generateProgram()), int64(0), int64(0), int64(0), make([]packet, 0)}
		computers = append(computers, &c)
	}

	return 0
}

func main() {

}
