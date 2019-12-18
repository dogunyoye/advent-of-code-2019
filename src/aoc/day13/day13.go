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

var numberOfOutputs = 0

type tileType int

const (
	empty            tileType = 0
	wall             tileType = 1
	block            tileType = 2
	horizontalPaddle tileType = 3
	ball             tileType = 4
)

type tile struct {
	X          int64
	Y          int64
	typeOfTile tileType
}

var part1 = false

var currentX = int64(0)
var currentY = int64(0)
var typeOfTile = empty

var tiles = make([]tile, 0)

var currentScore = int64(0)

var gameBall = tile{0, 0, ball}
var gamePaddle = tile{0, 0, horizontalPaddle}

var blocks = 0
var scoreOutput = 0

func copyArray(array []int64) []int64 {
	arrNew := make([]int64, 0)
	arrNew = append(arrNew, array...)

	// extend program's memory
	space := make([]int64, len(array)*1000)
	arrNew = append(arrNew, space...)
	return arrNew
}

func constructTile(output int64) {
	numberOfOutputs++

	switch numberOfOutputs {
	case 1:
		currentX = output
		break
	case 2:
		currentY = output
		break
	case 3:
		switch output {
		case 0:
			typeOfTile = empty
		case 1:
			typeOfTile = wall
		case 2:
			typeOfTile = block
		case 3:
			typeOfTile = horizontalPaddle
		case 4:
			typeOfTile = ball
		default:
			// current score output
			currentScore = output
		}

		var t = tile{currentX, currentY, typeOfTile}
		if part1 {
			tiles = append(tiles, t)
		} else {
			// do stuff
			if t.typeOfTile == ball {
				gameBall.X = t.X
				gameBall.Y = t.Y
			} else if t.typeOfTile == horizontalPaddle {
				gamePaddle.X = t.X
				gamePaddle.Y = t.Y
			}
		}

		numberOfOutputs = 0
		break
	}
}

func generateProgram() []int64 {
	file, err := os.Open("../../data/day13.txt")

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

func moveJoystick() int64 {
	// move left
	if gameBall.X < gamePaddle.X {
		return -1
	} else if gameBall.X > gamePaddle.X {
		return 1
	}

	return 0
}

func runOpcodeForParameterMode(opcode int64, opcodeIndex int64, program []int64, input int64) int {
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
		if !part1 {
			//game started
			input = moveJoystick()
		}

		if paramMode0 == 0 { // position mode
			program[program[opcodeIndex]] = input
		} else if paramMode0 == 1 { // immediate mode
			program[opcodeIndex] = input
		} else if paramMode0 == 2 { // relative mode
			program[program[opcodeIndex]+relativeBase] = input
		}
		return 2
	}

	if paramMode0 == 0 { // position mode
		firstOperand = program[program[opcodeIndex]]
	} else if paramMode0 == 1 { // immediate mode
		firstOperand = program[opcodeIndex]
	} else if paramMode0 == 2 { // relative mode
		firstOperand = program[program[opcodeIndex]+relativeBase]
	} else {
		fmt.Println("Unknown first param mode:", paramMode0)
	}

	// output instruction
	if opcode == 4 {
		output = firstOperand
		constructTile(output)
		return 2
	}

	// relative base adjustment instruction
	if opcode == 9 {
		relativeBase += firstOperand
		return 2
	}

	opcodeIndex++

	if paramMode1 == 0 { // position mode
		secondOperand = program[program[opcodeIndex]]
	} else if paramMode1 == 1 { // immediate mode
		secondOperand = program[opcodeIndex]
	} else if paramMode1 == 2 { // relative mode
		secondOperand = program[program[opcodeIndex]+relativeBase]
	} else {
		fmt.Println("Unknown second param mode:", paramMode1)
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
	} else if paramMode2 == 2 {
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

func runDiagnosticProgram(program []int64, input int64) {
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
			if !part1 {
				input = moveJoystick()
			}
			program[program[memPointer+1]] = input
		case 4:
			output = program[program[memPointer+1]]
			constructTile(output)
		case 5:
			fallthrough
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			opcodeJump = runOpcode(opcode, memPointer, program)
			break
		case 9:
			relativeBase += program[program[memPointer+1]]
			break
		case 99:
			return
		default:
			opcodeJump = runOpcodeForParameterMode(opcode, memPointer, program, input)
		}

		memPointer += int64(opcodeJump)
		opcode = program[memPointer]
	}
}

func main() {

	var program1 = copyArray(generateProgram())
	var program2 = copyArray(generateProgram())

	part1 = true

	runDiagnosticProgram(program1, 1)
	memPointer = 0
	relativeBase = 0

	blockTiles := 0
	for _, t := range tiles {
		if t.typeOfTile == block {
			blockTiles++
		}
	}

	fmt.Println("Part1:", blockTiles)

	part1 = false
	program2[0] = 2

	runDiagnosticProgram(program2, 2)
	fmt.Println("Part2:", currentScore)
}
