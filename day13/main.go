package main

import (
	"fmt"
	"log"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

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

func constructTile(output int64) {
	numberOfOutputs++

	switch numberOfOutputs {
	case 1:
		currentX = output
	case 2:
		currentY = output
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
			switch t.typeOfTile {
			case ball:
				gameBall.X = t.X
				gameBall.Y = t.Y
			case horizontalPaddle:
				gamePaddle.X = t.X
				gamePaddle.Y = t.Y
			}
		}

		numberOfOutputs = 0
	}
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

func runMachine(m *intcode.Machine, input int64) {
	m.In = func() int64 {
		if !part1 {
			input = moveJoystick()
		}
		return input
	}
	m.Out = func(v int64) { constructTile(v) }
	m.Run()
}

func main() {

	program, err := intcode.LoadProgram(puzzle.Path("day13.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	part1 = true

	runMachine(intcode.New(program), 1)

	blockTiles := 0
	for _, t := range tiles {
		if t.typeOfTile == block {
			blockTiles++
		}
	}

	fmt.Println("Part1:", blockTiles)

	part1 = false
	freePlay := intcode.New(program)
	freePlay.SetMem(0, 2)

	runMachine(freePlay, 2)
	fmt.Println("Part2:", currentScore)
}
