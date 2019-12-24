package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

type position struct {
	X int
	Y int
}

var eris = make(map[position]string)
var previousErisStates = make(map[string]struct{})

var terrain = [5][5]string{}

func parseAdjacentCells(currentInhabitant string, adjacents []position) string {
	var result = currentInhabitant

	var bugCount = 0
	for _, cell := range adjacents {
		_, exists := eris[cell]
		if exists && eris[cell] == "#" {
			bugCount++
		}
	}

	switch currentInhabitant {
	case ".":
		if bugCount == 1 || bugCount == 2 {
			result = "#"
		}
	case "#":
		if bugCount != 1 {
			result = "."
		}
	}

	return result
}

func calculateBiodiversity() int {
	var nextState = ""

	for {
		nextState = ""
		for k, v := range eris {
			var upAdj = position{k.X, k.Y - 1}
			var rightAdj = position{k.X + 1, k.Y}
			var bottomAdj = position{k.X, k.Y + 1}
			var leftAdj = position{k.X - 1, k.Y}

			var adjacents = []position{upAdj, rightAdj, bottomAdj, leftAdj}
			result := parseAdjacentCells(v, adjacents)

			terrain[k.X][k.Y] = result
		}

		for x := 0; x < 5; x++ {
			for y := 0; y < 5; y++ {
				nextState += terrain[x][y]
				eris[position{x, y}] = terrain[x][y]
			}
		}

		_, exists := previousErisStates[nextState]
		if !exists {
			previousErisStates[nextState] = struct{}{}
			continue
		}

		break
	}

	biodiversity := 0

	for i := 0; i < len(nextState); i++ {
		if nextState[i] == '#' {
			biodiversity += int(math.Pow(float64(2), float64(i)))
		}
	}

	return biodiversity
}

func main() {
	file, err := os.Open("../../data/day24.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	row := 0
	initialState := ""

	for scanner.Scan() {
		line := scanner.Text()
		for i := 0; i < len(line); i++ {

			var pos = position{row, i}
			eris[pos] = string(line[i])

			initialState += string(line[i])
		}

		row++
	}

	file.Close()

	previousErisStates[initialState] = struct{}{}

	fmt.Println("Part1:", calculateBiodiversity())
}
