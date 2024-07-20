package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"
)

type position struct {
	i int
	j int
}

type state struct {
	pos   position
	steps int
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func printMap(arena map[position]rune, depth int, width int) {
	for i := 0; i < depth; i++ {
		line := ""
		for j := 0; j < width; j++ {
			line += string(arena[position{i, j}])
		}
		fmt.Println(line)
	}
}

func findPortals(arena map[position]rune) (map[position]position, position, position) {
	portals := make(map[string][]position)
	for k, v := range arena {
		if v == '.' {
			var portalId = ""
			north := arena[position{k.i - 1, k.j}]
			east := arena[position{k.i, k.j + 1}]
			south := arena[position{k.i + 1, k.j}]
			west := arena[position{k.i, k.j - 1}]

			if unicode.IsLetter(north) {
				portalId = string(arena[position{k.i - 2, k.j}]) + string(north)
			} else if unicode.IsLetter(east) {
				portalId = string(east) + string(arena[position{k.i, k.j + 2}])
			} else if unicode.IsLetter(south) {
				portalId = string(south) + string(arena[position{k.i + 2, k.j}])
			} else if unicode.IsLetter(west) {
				portalId = string(arena[position{k.i, k.j - 2}]) + string(west)
			}

			if len(portalId) == 2 {
				_, exists := portals[portalId]
				if exists {
					portals[portalId] = append(portals[portalId], k)
				} else {
					positions := make([]position, 0)
					positions = append(positions, k)
					portals[portalId] = positions
				}
			}
		}
	}

	edges := make(map[position]position)
	var start = position{-1, -1}
	var end = position{-1, -1}

	for k, v := range portals {
		if k == "AA" {
			start = v[0]
		} else if k == "ZZ" {
			end = v[0]
		} else {
			edges[v[0]] = v[1]
			edges[v[1]] = v[0]
		}
	}

	return edges, start, end
}

func buildMap() (map[position]rune, int, int) {
	file, err := os.Open("../../data/day20.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var i = 0
	var width = 0
	arena := make(map[position]rune)

	for scanner.Scan() {
		for j, c := range scanner.Text() {
			arena[position{i, j}] = c
			width = j
		}

		i += 1
	}

	file.Close()
	return arena, i, width
}

func findFewestStepsFromAAToZZ() int {
	arena, _, _ := buildMap()
	portals, start, end := findPortals(arena)

	queue := make([]state, 0)
	visited := map[position]struct{}{}

	queue = append(queue, state{start, 0})
	visited[start] = struct{}{}

	for len(queue) != 0 {
		currentState := queue[0]
		queue = queue[1:]

		if currentState.pos == end {
			return currentState.steps
		}

		neighbours := [4]position{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

		for _, n := range neighbours {
			nextPosition := position{currentState.pos.i + n.i, currentState.pos.j + n.j}
			_, seen := visited[nextPosition]
			if !seen && arena[nextPosition] == '.' {
				visited[nextPosition] = struct{}{}
				teleportedPos, exists := portals[nextPosition]
				if exists {
					queue = append(queue, state{teleportedPos, currentState.steps + 2})
				} else {
					queue = append(queue, state{nextPosition, currentState.steps + 1})
				}
			}
		}
	}

	panic("No solution found!")
}

func main() {
	fmt.Println("Part1:", findFewestStepsFromAAToZZ())
}
