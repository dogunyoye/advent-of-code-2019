package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

type direction int

const (
	North direction = iota
	East
	South
	West
)

type position struct {
	i int
	j int
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func printLayers(layers map[int]map[position]rune) {
	keys := make([]int, 0)
	for k := range layers {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, l := range keys {
		fmt.Println("Layer", l)
		printMap(layers[l])
		fmt.Println()
	}
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func printMap(grid map[position]rune) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			fmt.Printf("%c", grid[position{i, j}])
		}
		fmt.Println()
	}
}

func hasBugs(grid map[position]rune) bool {
	for _, v := range grid {
		if v == '#' {
			return true
		}
	}

	return false
}

func isPositionBorder(pos position) bool {
	return pos.i == 0 || pos.j == 4 || pos.i == 4 || pos.j == 0
}

func getMinMaxLayer(layers map[int]map[position]rune) (int, int) {
	keys := make([]int, 0)
	for k := range layers {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	return keys[0], keys[len(keys)-1]
}

func adjacentToMiddle(grid map[position]rune) map[direction]rune {
	neighbours := make(map[direction]rune)
	neighbours[North] = grid[position{1, 2}]
	neighbours[East] = grid[position{2, 3}]
	neighbours[South] = grid[position{3, 2}]
	neighbours[West] = grid[position{2, 1}]
	return neighbours
}

func borders(grid map[position]rune) map[direction][]rune {
	borders := make(map[direction][]rune)
	borders[North] = make([]rune, 0)
	borders[East] = make([]rune, 0)
	borders[South] = make([]rune, 0)
	borders[West] = make([]rune, 0)

	for k, v := range grid {
		if k.i == 0 {
			borders[North] = append(borders[North], v)
		}

		if k.j == 4 {
			borders[East] = append(borders[East], v)
		}

		if k.i == 4 {
			borders[South] = append(borders[South], v)
		}

		if k.j == 0 {
			borders[West] = append(borders[West], v)
		}
	}

	return borders
}

func buildEmptyGrid() map[position]rune {
	grid := make(map[position]rune)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			val := '.'
			if i == 2 && j == 2 {
				val = '?'
			}

			grid[position{i, j}] = val
		}
	}

	return grid
}

func buildInitialState() (map[position]rune, map[string]struct{}) {
	file, err := os.Open("../../data/day24.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var eris = make(map[position]rune)
	var previousErisStates = make(map[string]struct{})

	i := 0
	initialState := ""

	for scanner.Scan() {
		line := scanner.Text()
		for j := 0; j < len(line); j++ {
			eris[position{i, j}] = rune(line[j])
			initialState += string(line[j])
		}
		i++
	}

	file.Close()

	previousErisStates[initialState] = struct{}{}
	return eris, previousErisStates
}

func parseAdjacentCells(currentInhabitant rune, adjacents []rune) rune {
	var result = currentInhabitant

	var bugCount = 0
	for _, cell := range adjacents {
		if cell == '#' {
			bugCount++
		}
	}

	switch currentInhabitant {
	case '.':
		if bugCount == 1 || bugCount == 2 {
			result = '#'
		}
	case '#':
		if bugCount != 1 {
			result = '.'
		}
	}

	return result
}

func calculateBiodiversity() int {
	eris, previousErisStates := buildInitialState()
	terrain := [5][5]rune{}

	var nextState = ""
	for {
		nextState = ""
		for k, v := range eris {
			var west = position{k.i, k.j - 1}
			var south = position{k.i + 1, k.j}
			var east = position{k.i, k.j + 1}
			var north = position{k.i - 1, k.j}

			var adjacents = []rune{eris[west], eris[south], eris[east], eris[north]}
			result := parseAdjacentCells(v, adjacents)

			terrain[k.i][k.j] = result
		}

		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				nextState += string(terrain[i][j])
				eris[position{i, j}] = terrain[i][j]
			}
		}

		_, exists := previousErisStates[nextState]
		if exists {
			break
		}

		previousErisStates[nextState] = struct{}{}
	}

	biodiversity := 0
	for i := 0; i < len(nextState); i++ {
		if nextState[i] == '#' {
			biodiversity += int(math.Pow(float64(2), float64(i)))
		}
	}

	return biodiversity
}

func findNumberOfBugsAfter200Minutes() int {
	layers := make(map[int]map[position]rune)
	initial, _ := buildInitialState()
	initial[position{2, 2}] = '?'
	minutes := 200

	layers[-1] = buildEmptyGrid()
	layers[0] = initial
	layers[1] = buildEmptyGrid()

	for minutes != 0 {
		// potentially expand the layers in each direction
		// generate new empty layer(s) above/below
		// the current min/max layer, if they have bugs
		min, max := getMinMaxLayer(layers)
		if hasBugs(layers[min]) {
			layers[min-1] = buildEmptyGrid()
		}

		if hasBugs(layers[max]) {
			layers[max+1] = buildEmptyGrid()
		}

		nextLayers := make(map[int]map[position]rune)

		for layer, grid := range layers {
			_, outsideExists := layers[layer-1]
			_, insideExists := layers[layer+1]

			nextGrid := make(map[position]rune)

			for k, v := range grid {

				if v == '?' {
					nextGrid[k] = '?'
					continue
				}

				neighbours := make([]rune, 0)

				north, northExists := grid[position{k.i - 1, k.j}]
				east, eastExists := grid[position{k.i, k.j + 1}]
				south, southExists := grid[position{k.i + 1, k.j}]
				west, westExists := grid[position{k.i, k.j - 1}]

				if outsideExists && isPositionBorder(k) {
					outside := adjacentToMiddle(layers[layer-1])

					if !northExists && !westExists { // north-west corner
						neighbours = append(neighbours, outside[North], east, south, outside[West])
					} else if !northExists && !eastExists { // north-east corner
						neighbours = append(neighbours, outside[North], outside[East], south, west)
					} else if !northExists && eastExists && southExists && westExists { // north
						neighbours = append(neighbours, outside[North], east, south, west)
					} else if northExists && !eastExists && southExists && westExists { // east
						neighbours = append(neighbours, north, outside[East], south, west)
					} else if !eastExists && !southExists { // south-east corner
						neighbours = append(neighbours, north, outside[East], outside[South], west)
					} else if northExists && eastExists && !southExists && westExists { // south
						neighbours = append(neighbours, north, east, outside[South], west)
					} else if !westExists && !southExists { // south-west corner
						neighbours = append(neighbours, north, east, outside[South], outside[West])
					} else if northExists && eastExists && southExists && !westExists { // west
						neighbours = append(neighbours, north, east, south, outside[West])
					}
				} else if insideExists && !isPositionBorder(k) {
					inside := borders(layers[layer+1])

					if north == '?' {
						neighbours = append(neighbours, inside[South]...)
						neighbours = append(neighbours, east, south, west)
					} else if east == '?' {
						neighbours = append(neighbours, inside[West]...)
						neighbours = append(neighbours, north, south, west)
					} else if south == '?' {
						neighbours = append(neighbours, inside[North]...)
						neighbours = append(neighbours, north, east, west)
					} else if west == '?' {
						neighbours = append(neighbours, inside[East]...)
						neighbours = append(neighbours, north, east, south)
					} else {
						neighbours = append(neighbours, north, east, south, west)
					}
				}

				nextGrid[k] = parseAdjacentCells(v, neighbours)
			}

			nextLayers[layer] = nextGrid
		}

		for k, v := range nextLayers {
			layers[k] = v
		}

		minutes -= 1
	}

	bugs := 0
	for _, grid := range layers {
		for _, v := range grid {
			if v == '#' {
				bugs += 1
			}
		}
	}

	return bugs
}

func main() {
	fmt.Println("Part1:", calculateBiodiversity())
	fmt.Println("Part2:", findNumberOfBugsAfter200Minutes())
}
