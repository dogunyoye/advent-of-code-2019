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

var memPointer = int64(0)
var output = int64(0)
var relativeBase = int64(0)

type position struct {
	X int
	Y int
}

var currentPos = position{-1, -1}

var positions = []position{}
var positionIndex = 0
var numOfInputs = 0

var pulled = 0

var part2 = false

var pullPointsList = []position{}
var pulledPointsMap = make(map[position]struct{})

func copyArray(array []int64) []int64 {
	arrNew := make([]int64, 0)
	arrNew = append(arrNew, array...)

	// extend program's memory
	space := make([]int64, len(array))
	arrNew = append(arrNew, space...)
	return arrNew
}

func generateProgram() []int64 {
	file, err := os.Open("../../data/day19.txt")

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

func setInput() int64 {
	currentPos = positions[positionIndex]
	var input = int64(0)
	switch numOfInputs {
	case 0:
		input = int64(currentPos.X)
		numOfInputs++
	case 1:
		input = int64(currentPos.Y)
		numOfInputs = 0
		positionIndex++
	}

	return input
}

func runOpcodeForParameterMode(opcode int64, opcodeIndex int64, program []int64) int {
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

		if output == 1 {
			pulled++
			if part2 {
				pos := position{currentPos.X, currentPos.Y}
				pullPointsList = append(pullPointsList, pos)
				pulledPointsMap[pos] = struct{}{}
			}
		}

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

func runDiagnosticProgram(program []int64) {
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
			program[program[memPointer+1]] = setInput()
		case 4:
			output = program[program[memPointer+1]]

			if output == 1 {
				pulled++
				if part2 {
					pos := position{currentPos.X, currentPos.Y}
					pullPointsList = append(pullPointsList, pos)
					pulledPointsMap[pos] = struct{}{}
				}
			}
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
			opcodeJump = runOpcodeForParameterMode(opcode, memPointer, program)
		}

		memPointer += int64(opcodeJump)
		opcode = program[memPointer]
	}
}

func main() {

	for x := 0; x < 50; x++ {
		for y := 0; y < 50; y++ {
			positions = append(positions, position{x, y})
		}
	}

	for i := 0; i < len(positions); i++ {
		var p = copyArray(generateProgram())
		runDiagnosticProgram(p)
		memPointer = 0
		relativeBase = 0
	}

	fmt.Println("Part1:", pulled)

	start := time.Now()

	// get ready for part 2
	positions = nil
	positions = []position{}
	numOfInputs = 0
	positionIndex = 0
	part2 = true

	for x := 0; x < 2000; x++ {
		for y := 0; y < 2000; y++ {
			positions = append(positions, position{x, y})
		}
	}

	program1 := generateProgram()

	// very slow loop, 4M iterations!
	// ~3mins completion time
	// todo: refactor solution
	for i := 0; i < len(positions); i++ {
		p := copyArray(program1)
		runDiagnosticProgram(p)
		memPointer = 0
		relativeBase = 0
	}

	pos := position{-1, -1}

	for _, p := range pullPointsList {

		_, alongX := pulledPointsMap[position{p.X + 99, p.Y}]
		_, alongY := pulledPointsMap[position{p.X, p.Y + 99}]

		if !alongX || !alongY {
			continue
		}

		pos = p
		break
	}

	fmt.Println("Part2:", (10000*pos.X)+pos.Y)

	elapsed := time.Since(start)
	fmt.Println("Part2 took:", elapsed)
}
