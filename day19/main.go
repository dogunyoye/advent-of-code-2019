package main

import (
	"fmt"
	"log"
	"time"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

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

func handleOutput(v int64) {
	if v == 1 {
		pulled++
		if part2 {
			pos := position{currentPos.X, currentPos.Y}
			pullPointsList = append(pullPointsList, pos)
			pulledPointsMap[pos] = struct{}{}
		}
	}
}

func probe(program []int64) {
	m := intcode.New(program)
	m.In = func() int64 { return setInput() }
	m.Out = handleOutput
	m.Run()
}

func main() {

	for x := 0; x < 50; x++ {
		for y := 0; y < 50; y++ {
			positions = append(positions, position{x, y})
		}
	}

	program, err := intcode.LoadProgram(puzzle.Path("day19.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	for i := 0; i < len(positions); i++ {
		probe(program)
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

	// very slow loop, 4M iterations!
	// ~3mins completion time
	// todo: refactor solution
	for i := 0; i < len(positions); i++ {
		probe(program)
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
