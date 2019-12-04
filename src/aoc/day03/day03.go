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

type position struct {
	X int
	Y int
}

type wire struct {
	currentPosition position
	instructions    []string
	traversed       []position
}

func mapWire(w wire) []position {
	for _, i := range w.instructions {
		direction := i[0:1]
		units, err := strconv.Atoi(i[1:])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		for m := 1; m <= units; m++ {
			switch direction {
			case "R":
				w.currentPosition.X++
			case "L":
				w.currentPosition.X--
			case "U":
				w.currentPosition.Y++
			case "D":
				w.currentPosition.Y--
			}

			w.traversed = append(w.traversed, position{w.currentPosition.X, w.currentPosition.Y})
		}
	}

	return w.traversed
}

func countStepsToIntersection(w wire, intersection position) int {
	steps := 0
	for _, i := range w.instructions {
		direction := i[0:1]
		units, err := strconv.Atoi(i[1:])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		for m := 1; m <= units; m++ {
			switch direction {
			case "R":
				w.currentPosition.X++
			case "L":
				w.currentPosition.X--
			case "U":
				w.currentPosition.Y++
			case "D":
				w.currentPosition.Y--
			}

			steps++
			if (w.currentPosition.X == intersection.X) && (w.currentPosition.Y == intersection.Y) {
				return steps
			}
		}

	}

	return steps
}

func main() {
	file, err := os.Open("../../data/day03.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var wire1 = wire{position{0, 0}, []string{}, []position{}}
	var wire2 = wire{position{0, 0}, []string{}, []position{}}

	scanner.Scan()
	wire1.instructions = strings.Split(scanner.Text(), ",")

	scanner.Scan()
	wire2.instructions = strings.Split(scanner.Text(), ",")

	wire1.traversed = mapWire(wire1)
	wire2.traversed = mapWire(wire2)

	traversedMap := make(map[position]bool)
	manhattanDistance := 0

	var intersections = []position{}
	fewestSteps := 0

	for _, t1 := range wire1.traversed {
		traversedMap[t1] = true
	}

	for _, t2 := range wire2.traversed {
		_, ok := traversedMap[t2]
		if ok {
			intersections = append(intersections, t2)
			calculatedDistance := math.Abs(float64(t2.X)) + math.Abs(float64(t2.Y))
			if manhattanDistance == 0 {
				manhattanDistance = int(calculatedDistance)
			} else {
				manhattanDistance = int(math.Min(calculatedDistance, float64(manhattanDistance)))
			}
		}
	}

	fmt.Println("Part1:", manhattanDistance)

	for _, i := range intersections {
		steps := countStepsToIntersection(wire1, i) + countStepsToIntersection(wire2, i)
		if fewestSteps == 0 {
			fewestSteps = steps
		} else {
			fewestSteps = int(math.Min(float64(steps), float64(fewestSteps)))
		}
	}

	fmt.Println("Part2:", fewestSteps)

	file.Close()
}
