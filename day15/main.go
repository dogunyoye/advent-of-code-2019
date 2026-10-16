package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

var spaceship = make(map[position]string)

var oxygenPos = position{-1, -1}

// currentInput holds the movement command for the in-progress step.
// The machine may consume it more than once per step, exactly like the
// original per-call input parameter.
var currentInput int64

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

// stepMachine feeds one movement command and runs until the next
// output, mirroring the original pause-before-input protocol.
func stepMachine(m *intcode.Machine, dir direction) int64 {
	currentInput = int64(dir)
	res, halted := m.RunUntilOutput()
	if halted {
		fmt.Println("Should not halt")
		os.Exit(2)
	}
	return res
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

func moveInDirection(m *intcode.Machine, currentPos position, dir direction) {
	switch dir {
	case north:
		exploreShip(m, position{currentPos.X - 1, currentPos.Y})
	case east:
		exploreShip(m, position{currentPos.X, currentPos.Y + 1})
	case south:
		exploreShip(m, position{currentPos.X + 1, currentPos.Y})
	case west:
		exploreShip(m, position{currentPos.X, currentPos.Y - 1})
	}

	stepMachine(m, opposite(dir))
}

func exploreShip(m *intcode.Machine, currentPos position) {
	_, exists := spaceship[currentPos]
	if exists {
		return
	}

	for _, dir := range directions {
		res := stepMachine(m, dir)

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
			moveInDirection(m, currentPos, dir)

		// moved into space and found oxygen (O)
		case 2:
			oxygenPos = currentPos
			spaceship[position{currentPos.X, currentPos.Y}] = "O"
			moveInDirection(m, currentPos, dir)

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

		if val, exists := spaceship[pos]; exists && (val != "#") {
			neighbours = append(neighbours, pos)
		}
	}

	return neighbours
}

func calculatePath(lastCheckedNode position, previous map[position]position) []position {
	var path = []position{}
	var temp = lastCheckedNode

	path = append(path, temp)
	for {
		if prev, exists := previous[temp]; exists {
			path = append(path, prev)
			temp = prev
		} else {
			return path
		}
	}
}

func aStarSearch(startPos position, endPos position) []position {
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
	return []position{}
}

func fillWithOxygen(oxygenPos position) int {
	open := []position{}
	closed := make(map[position]struct{})

	open = append(open, oxygenPos)
	var minutes = 0

	for len(open) > 0 {
		var toRemove = []position{}
		for _, p := range open {
			toRemove = append(toRemove, p)
			currentNeighbours := neighbours(p)
			for _, n := range currentNeighbours {
				if _, exists := closed[n]; !exists {
					open = append(open, n)
				}
			}
			closed[p] = struct{}{}
		}
		minutes++

		for _, r := range toRemove {
			open = remove(open, r)
		}
	}

	return minutes
}

func main() {

	program, err := intcode.LoadProgram(puzzle.Path("day15.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	m := intcode.New(program)
	m.In = func() int64 { return currentInput }

	var startPos = position{50, 50}

	exploreShip(m, startPos)

	path := aStarSearch(startPos, oxygenPos)

	fmt.Println("Part1:", len(path))
	fmt.Println("Part2:", fillWithOxygen(oxygenPos))
}
