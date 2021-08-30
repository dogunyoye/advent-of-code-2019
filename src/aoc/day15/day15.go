package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var memPointer = int64(0)
var output = int64(0)
var relativeBase = int64(0)

var spaceship = make(map[position]string)

var oxygenPos = position{-1, -1}

type direction int

const (
	north direction = 1
	south direction = 2
	west  direction = 3
	east  direction = 4
)

type position struct {
	X int
	Y int
}

var directions = []direction{
	north,
	east,
	south,
	west,
}

func copyArray(array []int64) []int64 {
	arrNew := make([]int64, 0)
	arrNew = append(arrNew, array...)
	return arrNew
}

func generateProgram() []int64 {
	file, err := os.Open("../../data/day15.txt")

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

func runDiagnosticProgram(program []int64, input int64) int64 {
	var opcode = program[memPointer]
	var hasOutput = false

	for {
		var opcodeJump = 2

		switch opcode {
		case 1:
			fallthrough
		case 2:
			opcodeJump = runOpcode(opcode, memPointer, program)
			break
		case 3:
			if hasOutput {
				return output
			}
			program[program[memPointer+1]] = input
		case 4:
			output = program[program[memPointer+1]]
			hasOutput = true
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
			fmt.Println("Should not halt")
			os.Exit(2)
		default:
			opcodeJump = runOpcodeForParameterMode(opcode, memPointer, program, input)
		}

		memPointer += int64(opcodeJump)
		opcode = program[memPointer]
	}
}

func opposite(dir direction) direction {
	switch dir {
	case north:
		return south
	case east:
		return west
	case south:
		return north
	case west:
		return east
	}

	fmt.Println("Invalid output")
	return -1
}

func moveInDirection(program []int64, currentPos position, dir direction) {
	switch dir {
	case north:
		exploreShip(program, position{currentPos.X - 1, currentPos.Y})
	case east:
		exploreShip(program, position{currentPos.X, currentPos.Y + 1})
	case south:
		exploreShip(program, position{currentPos.X + 1, currentPos.Y})
	case west:
		exploreShip(program, position{currentPos.X, currentPos.Y - 1})
	}

	runDiagnosticProgram(program, int64(opposite(dir)))
}

func exploreShip(program []int64, currentPos position) {
	_, exists := spaceship[currentPos]
	if exists {
		return
	}

	for _, dir := range directions {
		res := runDiagnosticProgram(program, int64(dir))

		switch res {

		// hit wall (#)
		case 0:
			var pos = position{-1, -1}
			switch dir {
			case north:
				pos = position{currentPos.X - 1, currentPos.Y}
			case east:
				pos = position{currentPos.X, currentPos.Y + 1}
			case south:
				pos = position{currentPos.X + 1, currentPos.Y}
			case west:
				pos = position{currentPos.X, currentPos.Y - 1}
			}

			spaceship[pos] = "#"

		// moved into space (.)
		case 1:
			spaceship[position{currentPos.X, currentPos.Y}] = "."
			moveInDirection(program, currentPos, dir)

		// moved into space and found oxygen (O)
		case 2:
			oxygenPos = currentPos
			spaceship[position{currentPos.X, currentPos.Y}] = "O"
			moveInDirection(program, currentPos, dir)

		default:
			fmt.Println("Invalid output")
			os.Exit(2)
		}
	}
}

func remove(l []position, item position) []position {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func includes(positions []position, toCheck position) bool {
	for _, p := range positions {
		if p == toCheck {
			return true
		}
	}

	return false
}

func heuristic(neighbour position, current position) int {
	return int(math.Abs(float64(neighbour.X-current.X)) + math.Abs(float64(neighbour.Y-current.Y)))
}

func neighbours(currentPos position) []position {
	var neighbours = []position{}

	for _, dir := range directions {

		var pos = position{-1, -1}
		switch dir {
		case north:
			pos = position{currentPos.X - 1, currentPos.Y}
		case east:
			pos = position{currentPos.X, currentPos.Y + 1}
		case south:
			pos = position{currentPos.X + 1, currentPos.Y}
		case west:
			pos = position{currentPos.X, currentPos.Y - 1}
		}

		val, exists := spaceship[pos]
		if exists && (val != "#") {
			neighbours = append(neighbours, pos)
		}
	}

	return neighbours
}

func calculatePath(lastCheckedNode position, previous map[position]position) int {
	var path = []position{}
	var temp = lastCheckedNode

	path = append(path, temp)
	for {
		if prev, exists := previous[temp]; exists {
			path = append(path, prev)
			temp = prev
		} else {
			return len(path)
		}
	}
}

func aStarSearch(startPos position, endPos position) int {
	var openSet = []position{}
	var closedSet = []position{}
	var lastCheckedNode = startPos

	hScore := make(map[position]int)
	gScore := make(map[position]int)
	fScore := make(map[position]int)

	previous := make(map[position]position)
	openSet = append(openSet, startPos)

	for len(openSet) > 0 {

		var winner = 0
		for i := 1; i < len(openSet); i++ {
			if fScore[openSet[i]] < fScore[openSet[winner]] {
				winner = i
			}

			if fScore[openSet[i]] == fScore[openSet[winner]] {
				if gScore[openSet[i]] > gScore[openSet[winner]] {
					winner = i
				}
			}
		}

		var current = openSet[winner]
		lastCheckedNode = current

		if current == endPos {
			// found solution
			return calculatePath(lastCheckedNode, previous)
		}

		openSet = remove(openSet, current)
		closedSet = append(closedSet, current)

		currNeighbours := neighbours(current)

		for i := 0; i < len(currNeighbours); i++ {
			var n = currNeighbours[i]

			if !includes(closedSet, n) {
				var tempG = gScore[current] + heuristic(n, current)

				if !includes(openSet, n) {
					openSet = append(openSet, n)
				} else if tempG >= gScore[n] {
					continue
				}

				gScore[n] = tempG
				hScore[n] = heuristic(n, endPos)
				fScore[n] = gScore[n] + hScore[n]

				previous[n] = current
			}
		}
	}

	// no solution!
	return -1
}

func main() {

	var program = copyArray(generateProgram())
	var startPos = position{50, 50}

	exploreShip(program, startPos)

	memPointer = 0
	relativeBase = 0

	fmt.Println("Part1:", aStarSearch(startPos, oxygenPos))
}
