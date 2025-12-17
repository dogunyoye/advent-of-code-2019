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

type recursiveState struct {
	pos     position
	steps   int
	arena   map[position]rune
	level   int
	visited map[levelState]struct{}
}

type levelState struct {
	pos   position
	level int
}

type portalPosition int

const (
	Unknown portalPosition = iota
	Inner
	Outer
)

type portal struct {
	portalPos portalPosition
	pos       position
	label     string
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

func findPortals(arena map[position]rune, depth int, width int) (map[portal]portal, position, position) {
	portals := make(map[string][]portal)
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

				var portalPos = Inner
				if k.i == 2 || k.i == depth-3 || k.j == 2 || k.j == width-2 {
					portalPos = Outer
				}

				_, exists := portals[portalId]
				if exists {
					portals[portalId] = append(portals[portalId], portal{portalPos, k, portalId})
				} else {
					positions := make([]portal, 0)
					positions = append(positions, portal{portalPos, k, portalId})
					portals[portalId] = positions
				}
			}
		}
	}

	edges := make(map[portal]portal)
	var start = position{-1, -1}
	var end = position{-1, -1}

	for k, v := range portals {
		switch k {
		case "AA":
			start = v[0].pos
		case "ZZ":
			end = v[0].pos
		default:
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

func isPortal(portals map[portal]portal, currentPos position) (portal, portal, bool) {
	for k, v := range portals {
		if k.pos == currentPos {
			return k, v, true
		}
	}

	return portal{Unknown, position{-1, -1}, "?"}, portal{Unknown, position{-1, -1}, "?"}, false
}

func copyArena(arena map[position]rune) map[position]rune {
	copy := make(map[position]rune)
	for k, v := range arena {
		copy[position{k.i, k.j}] = v
	}
	return copy
}

func levelZeroArena(arena map[position]rune, portals map[portal]portal, start position, end position) map[position]rune {
	copy := copyArena(arena)
	for k := range portals {
		if k.portalPos == Outer {
			copy[k.pos] = '#'
		}
	}
	for k := range copy {
		if k == start || k == end {
			copy[k] = '.'
		}
	}
	return copy
}

func aboveZeroArena(arena map[position]rune, portals map[portal]portal, start position, end position) map[position]rune {
	copy := copyArena(arena)
	for p := range portals {
		copy[p.pos] = '.'
	}

	for k := range copy {
		if k == start || k == end {
			copy[k] = '#'
		}
	}
	return copy
}

func findFewestStepsFromAAToZZ() int {
	arena, depth, width := buildMap()
	portals, start, end := findPortals(arena, depth, width)

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
				_, teleportedPortal, exists := isPortal(portals, nextPosition)
				if exists {
					queue = append(queue, state{teleportedPortal.pos, currentState.steps + 2})
				} else {
					queue = append(queue, state{nextPosition, currentState.steps + 1})
				}
			}
		}
	}

	panic("No solution found!")
}

// takes some time to complete
// not terrible, but could be optimised
func bfs(start position, initialArena map[position]rune, end position, portals map[portal]portal) int {

	visited := make(map[levelState]struct{})
	visited[levelState{start, 0}] = struct{}{}

	queue := make([]recursiveState, 0)
	queue = append(queue, recursiveState{start, 0, initialArena, 0, visited})

	aboveZero := aboveZeroArena(initialArena, portals, start, end)

	for len(queue) != 0 {
		currentState := queue[0]
		queue = queue[1:]

		if currentState.pos == end && currentState.level == 0 {
			return currentState.steps
		}

		neighbours := [4]position{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

		for _, n := range neighbours {
			nextPosition := position{currentState.pos.i + n.i, currentState.pos.j + n.j}
			_, seen := currentState.visited[levelState{nextPosition, currentState.level}]

			if !seen && currentState.arena[nextPosition] == '.' {
				currentState.visited[levelState{nextPosition, currentState.level}] = struct{}{}
				currentPortal, teleportedPortal, exists := isPortal(portals, nextPosition)

				if exists {
					var nextLevel = -1
					if currentPortal.portalPos == Inner {
						nextLevel = currentState.level + 1
					} else {
						nextLevel = currentState.level - 1
					}

					var nextArena = currentState.arena

					if nextLevel == 0 {
						nextArena = initialArena
					} else {
						nextArena = aboveZero
					}

					newVisited := make(map[levelState]struct{})
					newVisited[levelState{teleportedPortal.pos, nextLevel}] = struct{}{}
					queue = append(queue, recursiveState{teleportedPortal.pos, currentState.steps + 2, nextArena, nextLevel, newVisited})
				} else {
					queue = append(queue, recursiveState{nextPosition, currentState.steps + 1, currentState.arena, currentState.level, currentState.visited})
				}
			}
		}
	}

	panic("No solution found!")
}

func findFewestStepsFromAAToZZWithRecursiveTeleporting() int {
	arena, depth, width := buildMap()
	portals, start, end := findPortals(arena, depth, width)
	initialArena := levelZeroArena(arena, portals, start, end)

	return bfs(start, initialArena, end, portals)
}

func main() {
	fmt.Println("Part1:", findFewestStepsFromAAToZZ())
	fmt.Println("Part2:", findFewestStepsFromAAToZZWithRecursiveTeleporting())
}
