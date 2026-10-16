package main

import (
	"fmt"
	"log"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

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

func setRobot(output int64) {
	paintingRobot.outputs++

	if output > 1 {
		fmt.Println("??, output", output)
	}

	// on every second output set robot direction
	if paintingRobot.outputs%2 == 0 {
		switch paintingRobot.facing {
		case up:
			switch output {
			case 0:
				paintingRobot.facing = left
				paintingRobot.pos.X--
			case 1:
				paintingRobot.facing = right
				paintingRobot.pos.X++
			}
		case right:
			switch output {
			case 0:
				paintingRobot.facing = up
				paintingRobot.pos.Y++
			case 1:
				paintingRobot.facing = down
				paintingRobot.pos.Y--
			}
		case down:
			switch output {
			case 0:
				paintingRobot.facing = right
				paintingRobot.pos.X++
			case 1:
				paintingRobot.facing = left
				paintingRobot.pos.X--
			}
		case left:
			switch output {
			case 0:
				paintingRobot.facing = down
				paintingRobot.pos.Y--
			case 1:
				paintingRobot.facing = up
				paintingRobot.pos.Y++
			}
		}
	} else { // on every first output set the colour for the position

		var colour = ""
		switch output {
		case 0:
			colour = "."
		case 1:
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
		switch colour {
		case ".":
			input = 0
		case "#":
			input = 1
		}
	}

	return input
}

func runRobot(program []int64) {
	m := intcode.New(program)
	m.In = func() int64 { return int64(getInputFromRobot()) }
	m.Out = func(v int64) { setRobot(v) }
	m.Run()
}

func main() {

	program, err := intcode.LoadProgram(puzzle.Path("day11.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	runRobot(program)

	fmt.Println("Part1:", len(paintingRobot.traversed))

	// re-initialise the painting robot for part2
	paintingRobot = robot{up, 0, position{50, 50}, make(map[position]string)}
	// start on white square
	paintingRobot.traversed[paintingRobot.pos] = "#"

	runRobot(program)

	for k, v := range paintingRobot.traversed {
		visualise[k.Y][k.X] = v
	}

	fmt.Println("Part2:")

	for i := 0; i < 150; i++ {
		hasPrinted := false
		for j := 0; j < 150; j++ {
			if visualise[j][i] != "" {
				fmt.Print(visualise[j][i])
				hasPrinted = true
			}
		}

		if hasPrinted {
			fmt.Println()
		}
	}

}
