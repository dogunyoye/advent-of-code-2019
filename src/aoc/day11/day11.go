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

var visualise = [150][150]string{}

type direction int

const (
	up    direction = 0
	right direction = 1
	down  direction = 2
	left  direction = 3
)

type position struct {
	X int
	Y int
}

type robot struct {
	facing    direction
	outputs   int
	pos       position
	traversed map[position]string
}

// initialise robot
var paintingRobot = robot{up, 0, position{100, 100}, make(map[position]string)}

func copyArray(array []int64) []int64 {
	arrNew := make([]int64, 0)
	arrNew = append(arrNew, array...)

	// extend program's memory
	space := make([]int64, len(array)*1000)
	arrNew = append(arrNew, space...)
	return arrNew
}

func setRobot(output int64) {
	paintingRobot.outputs++

	if output > 1 {
		fmt.Println("??, output", output)
	}

	// on every second output set robot direction
	if paintingRobot.outputs%2 == 0 {
		switch paintingRobot.facing {
		case up:
			if output == 0 {
				paintingRobot.facing = left
				paintingRobot.pos.X--
			} else if output == 1 {
				paintingRobot.facing = right
				paintingRobot.pos.X++
			}
		case right:
			if output == 0 {
				paintingRobot.facing = up
				paintingRobot.pos.Y++
			} else if output == 1 {
				paintingRobot.facing = down
				paintingRobot.pos.Y--
			}
		case down:
			if output == 0 {
				paintingRobot.facing = right
				paintingRobot.pos.X++
			} else if output == 1 {
				paintingRobot.facing = left
				paintingRobot.pos.X--
			}
		case left:
			if output == 0 {
				paintingRobot.facing = down
				paintingRobot.pos.Y--
			} else if output == 1 {
				paintingRobot.facing = up
				paintingRobot.pos.Y++
			}
		}
	} else { // on every first output set the colour for the position

		var colour = ""
		if output == 0 {
			colour = "."
		} else if output == 1 {
			colour = "#"
		}

		_, exists := paintingRobot.traversed[paintingRobot.pos]
		if exists {
			delete(paintingRobot.traversed, paintingRobot.pos)
		}

		paintingRobot.traversed[paintingRobot.pos] = colour
	}
}

func getInputFromRobot() int {
	input := 0
	var colour = ""

	_, exists := paintingRobot.traversed[paintingRobot.pos]
	if exists {
		colour = paintingRobot.traversed[paintingRobot.pos]
		if colour == "." {
			input = 0
		} else if colour == "#" {
			input = 1
		}
	}

	return input
}

func generateProgram() []int64 {
	file, err := os.Open("../../data/day11.txt")

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

		input := int64(getInputFromRobot())

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
		setRobot(output)

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
			break
		case 3:
			input := int64(getInputFromRobot())
			program[program[memPointer+1]] = input
		case 4:
			output = program[program[memPointer+1]]
			setRobot(output)
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
			opcodeJump = runOpcodeForParameterMode(opcode, memPointer, program)
		}

		memPointer += int64(opcodeJump)
		opcode = program[memPointer]
	}
}

func main() {

	var program1 = copyArray(generateProgram())
	var program2 = copyArray(program1)

	runDiagnosticProgram(program1)
	memPointer = 0
	relativeBase = 0
	output = 0

	fmt.Println("Part1:", len(paintingRobot.traversed))

	// re-initialise the painting robot for part2
	paintingRobot = robot{up, 0, position{50, 50}, make(map[position]string)}
	// start on white square
	paintingRobot.traversed[paintingRobot.pos] = "#"

	runDiagnosticProgram(program2)

	for k, v := range paintingRobot.traversed {
		visualise[k.Y][k.X] = v
	}

	fmt.Println("Part2:")

	for i := 0; i < 150; i++ {
		for j := 0; j < 150; j++ {
			fmt.Print(visualise[j][i])
		}
		fmt.Println()
	}

}
