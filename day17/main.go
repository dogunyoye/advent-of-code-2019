package main

import (
	"fmt"
	"log"
	"strconv"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

type position struct {
	X int
	Y int
}

type direction int

const (
	up    direction = 0
	right direction = 1
	down  direction = 2
	left  direction = 3
)

type robot struct {
	facing    direction
	pos       position
	traversed map[position]struct{}
}

var view = []byte{}
var viewMap = make(map[position]string)

var currentX = 0
var currentY = 0

var scaffoldingCount = 0

var vacuumRobot = robot{up, position{-1, -1}, make(map[position]struct{})}

var inputIndex = 0
var programInput = []int{}

var part1 = false

func buildView(output int64) {
	if output == 10 {
		currentX = 0
		currentY++

		view = append(view, byte(output))
		return
	}

	view = append(view, byte(output))
	viewMap[position{currentX, currentY}] = string(byte(output))

	currentX++
}

func setRobotDirectionAndTraverse(instructions []string) []string {

	for len(vacuumRobot.traversed) != scaffoldingCount {
		e1, e2, e3, e4 := false, false, false, false
		moves := 0

		_, northExists := viewMap[position{vacuumRobot.pos.X, vacuumRobot.pos.Y - 1}]
		if northExists && (viewMap[position{vacuumRobot.pos.X, vacuumRobot.pos.Y - 1}] == "#") {
			e1 = true
		}

		_, southExists := viewMap[position{vacuumRobot.pos.X, vacuumRobot.pos.Y + 1}]
		if southExists && (viewMap[position{vacuumRobot.pos.X, vacuumRobot.pos.Y + 1}] == "#") {
			e2 = true
		}

		_, eastExists := viewMap[position{vacuumRobot.pos.X + 1, vacuumRobot.pos.Y}]
		if eastExists && (viewMap[position{vacuumRobot.pos.X + 1, vacuumRobot.pos.Y}] == "#") {
			e3 = true
		}

		_, westExists := viewMap[position{vacuumRobot.pos.X - 1, vacuumRobot.pos.Y}]
		if westExists && (viewMap[position{vacuumRobot.pos.X - 1, vacuumRobot.pos.Y}] == "#") {
			e4 = true
		}

		switch vacuumRobot.facing {

		case down:
			if e4 {
				vacuumRobot.facing = left
				vacuumRobot.pos.X--
				instructions = append(instructions, "R")
			} else if e3 {
				vacuumRobot.facing = right
				vacuumRobot.pos.X++
				instructions = append(instructions, "L")
			} else {
				vacuumRobot.pos.Y++
			}

		case up:
			if e4 {
				vacuumRobot.facing = left
				vacuumRobot.pos.X--
				instructions = append(instructions, "L")
			} else if e3 {
				vacuumRobot.facing = right
				vacuumRobot.pos.X++
				instructions = append(instructions, "R")
			} else {
				vacuumRobot.pos.Y--
			}

		case left:
			if e1 {
				vacuumRobot.facing = up
				vacuumRobot.pos.Y--
				instructions = append(instructions, "R")
			} else if e2 {
				vacuumRobot.facing = down
				vacuumRobot.pos.Y++
				instructions = append(instructions, "L")
			} else {
				vacuumRobot.pos.X--
			}

		case right:
			if e1 {
				vacuumRobot.facing = up
				vacuumRobot.pos.Y--
				instructions = append(instructions, "L")
			} else if e2 {
				vacuumRobot.facing = down
				vacuumRobot.pos.Y++
				instructions = append(instructions, "R")
			} else {
				vacuumRobot.pos.X++
			}
		}

		vacuumRobot.traversed[vacuumRobot.pos] = struct{}{}
		moves++

		scaffoldAhead := true

		for scaffoldAhead {
			var nextPos = vacuumRobot.pos
			switch vacuumRobot.facing {
			case up:
				nextPos.X = vacuumRobot.pos.X
				nextPos.Y = vacuumRobot.pos.Y - 1
			case down:
				nextPos.X = vacuumRobot.pos.X
				nextPos.Y = vacuumRobot.pos.Y + 1
			case left:
				nextPos.X = vacuumRobot.pos.X - 1
				nextPos.Y = vacuumRobot.pos.Y
			case right:
				nextPos.X = vacuumRobot.pos.X + 1
				nextPos.Y = vacuumRobot.pos.Y
			}

			_, exists := viewMap[nextPos]
			if exists && viewMap[nextPos] == "#" {
				vacuumRobot.traversed[nextPos] = struct{}{}
				vacuumRobot.pos = nextPos
				moves++
			} else {
				scaffoldAhead = false
			}
		}

		instructions = append(instructions, strconv.Itoa(moves))
	}

	return instructions
}

func runProgram(program []int64) int64 {
	var last int64
	m := intcode.New(program)
	m.In = func() int64 {
		input := int64(programInput[inputIndex])
		inputIndex++
		return input
	}
	m.Out = func(v int64) {
		if part1 {
			buildView(v)
		}
		last = v
	}
	m.Run()
	return last
}

func main() {

	part1 = true
	program, err := intcode.LoadProgram(puzzle.Path("day17.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	last := runProgram(program)

	total := 0

	for k, v := range viewMap {
		switch v {
		case "#":
			scaffoldingCount++
			e1, e2, e3, e4 := false, false, false, false

			_, northExists := viewMap[position{k.X, k.Y - 1}]
			if northExists && (viewMap[position{k.X, k.Y - 1}] == "#") {
				e1 = true
			}

			_, southExists := viewMap[position{k.X, k.Y + 1}]
			if southExists && (viewMap[position{k.X, k.Y + 1}] == "#") {
				e2 = true
			}

			_, eastExists := viewMap[position{k.X + 1, k.Y}]
			if eastExists && (viewMap[position{k.X + 1, k.Y}] == "#") {
				e3 = true
			}

			_, westExists := viewMap[position{k.X - 1, k.Y}]
			if westExists && (viewMap[position{k.X - 1, k.Y}] == "#") {
				e4 = true
			}

			if e1 && e2 && e3 && e4 {
				total += k.X * k.Y
			}
		case "^", ">", "V", "<":
			switch v {
			case "^":
				vacuumRobot.facing = up
			case ">":
				vacuumRobot.facing = right
			case "V":
				vacuumRobot.facing = down
			case "<":
				vacuumRobot.facing = left
			}

			vacuumRobot.pos.X = k.X
			vacuumRobot.pos.Y = k.Y
		}
	}

	fmt.Println("Part1:", total)
	part1 = false

	var instructions = []string{}
	instructions = setRobotDirectionAndTraverse(instructions)

	instructionString := ""

	for x, i := range instructions {
		instructionString += i
		if x != len(instructions)-1 {
			instructionString += ","
		}
	}

	// function instructions specific to my input
	var a = []byte("R,6,L,8,R,8\n")
	var b = []byte("R,4,R,6,R,6,R,4,R,4\n")
	var c = []byte("L,8,R,6,L,10,L,10\n")

	//A,A,B,C,B,C,B,C,A,C
	var mainRoutine = []byte("A,A,B,C,B,C,B,C,A,C\n")

	for _, i := range mainRoutine {
		programInput = append(programInput, int(i))
	}

	for _, i := range a {
		programInput = append(programInput, int(i))
	}

	for _, i := range b {
		programInput = append(programInput, int(i))
	}

	for _, i := range c {
		programInput = append(programInput, int(i))
	}

	programInput = append(programInput, 110)
	programInput = append(programInput, 10)

	m := intcode.New(program)
	m.SetMem(0, 2)
	m.In = func() int64 {
		input := int64(programInput[inputIndex])
		inputIndex++
		return input
	}
	m.Out = func(v int64) { last = v }
	m.Run()

	fmt.Println("Part2:", last)
}
