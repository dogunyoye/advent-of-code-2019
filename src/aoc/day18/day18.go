package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"unicode"
)

type position struct {
	i int
	j int
}

type positionState struct {
	keys  string
	pos   position
	steps int
}

type state struct {
	pos  position
	keys string
}

type quadrant struct {
	arena map[position]rune
	keys  string
}

func printMap(arena map[position]rune, depthStart int, depthEnd int, widthStart int, widthEnd int) {
	for i := depthStart; i < depthEnd; i++ {
		line := ""
		for j := widthStart; j < widthEnd; j++ {
			line += string(arena[position{i, j}])
		}
		fmt.Println(line)
	}
}

func buildMap() (map[position]rune, int, int, position) {
	file, err := os.Open("../../data/day18.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	arena := make(map[position]rune)
	var i = 0
	var width = 0
	var start = position{-1, -1}

	for scanner.Scan() {
		for j, c := range scanner.Text() {
			pos := position{i, j}
			arena[pos] = c
			width = j

			if c == '@' {
				start = pos
				arena[pos] = '.'
			}
		}
		i += 1
	}

	return arena, i, width + 1, start
}

func buildModifiedMap() (map[position]rune, int, int, []position) {
	arena, depth, width, start := buildMap()
	startPositions := make([]position, 0)

	arena[start] = '#'
	neighbours := [4]position{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}  // N, E, S, W
	diagonals := [4]position{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}} // NW, NE, SE, SW

	for _, n := range neighbours {
		p := position{start.i + n.i, start.j + n.j}
		arena[p] = '#'
	}

	for _, n := range diagonals {
		startPositions = append(startPositions, position{start.i + n.i, start.j + n.j})
	}

	return arena, depth, width, startPositions
}

func buildQuadrant(arena map[position]rune, depthStart int, depthEnd, widthStart int, widthEnd int) map[position]rune {
	quadrantArena := make(map[position]rune)

	for i := depthStart; i < depthEnd; i++ {
		for j := widthStart; j < widthEnd; j++ {
			p := position{i, j}
			quadrantArena[p] = arena[p]
		}
	}

	return quadrantArena
}

func isDoor(value rune) bool {
	return value >= 65 && value <= 90
}

func isDoorKey(value rune) bool {
	return value >= 97 && value <= 122
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func findKeys(arena map[position]rune) (int, string) {
	keys := 0
	heldKeys := strings.Repeat("0", 26)
	for _, v := range arena {
		if isDoorKey(v) {
			keys += 1
			heldKeys = replaceAtIndex(heldKeys, '1', f(v))
		}
	}
	return keys, heldKeys
}

func f(in rune) int {
	return int(in - 'a')
}

func obtainKey(keys string, pos int) string {
	return replaceAtIndex(keys, '1', pos)
}

func hasKey(keys string, k int) bool {
	return keys[k] == '1'
}

func collectKeys(arena map[position]rune, startState positionState) int {
	queue := make([]positionState, 0)
	visited := map[state]struct{}{}

	queue = append(queue, startState)
	visited[state{startState.pos, strings.Repeat("0", 26)}] = struct{}{}

	distance := math.MaxUint32
	_, allKeys := findKeys(arena)

	for len(queue) != 0 {

		currentPositionState := queue[0]
		queue = queue[1:]

		currentPosition := currentPositionState.pos
		currentSteps := currentPositionState.steps
		currentKeys := currentPositionState.keys

		if currentKeys == allKeys {
			if currentSteps < distance {
				distance = currentSteps
			}
			continue
		}

		neighbours := [4]position{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

		for _, n := range neighbours {
			nextPosition := position{currentPosition.i + n.i, currentPosition.j + n.j}
			nextKeys := currentKeys
			nextState := state{nextPosition, nextKeys}

			value := arena[nextPosition]
			if value == '#' {
				continue
			}

			_, exists := visited[nextState]
			if !exists {
				if isDoorKey(value) {
					nextKeys = obtainKey(nextKeys, f(value))
					queue = append(queue, positionState{nextKeys, nextPosition, currentSteps + 1})
					visited[state{nextPosition, nextKeys}] = struct{}{}
				} else {
					if (isDoor(value) && hasKey(nextKeys, f(unicode.ToLower(value)))) || value == '.' {
						queue = append(queue, positionState{nextKeys, nextPosition, currentSteps + 1})
						visited[state{nextPosition, nextKeys}] = struct{}{}
					}
				}
			}
		}
	}

	return distance
}

func findShortestPathToCollectAllKeys() int {
	arena, _, _, start := buildMap()
	keys := strings.Repeat("0", 26)
	result := collectKeys(arena, positionState{keys, start, 0})
	return result
}

func findShortestPathToCollectAllKeysWithFourRobots() int {
	arena, depth, width, startPositions := buildModifiedMap()
	printMap(arena, 0, depth, 0, width)
	fmt.Println(startPositions)

	//quadrants := make([]map[position]rune, 0)
	halfDepth := int(math.Round(float64(depth) / 2))
	halfWidth := int(math.Round(float64(width) / 2))

	fmt.Println(halfDepth)
	fmt.Println(halfWidth)

	quadrantDimensions := [][]int{
		{0, halfDepth, 0, halfWidth},
		{0, halfDepth, halfWidth - 1, width},
		{halfDepth - 1, depth, halfWidth - 1, width},
		{halfDepth - 1, depth, 0, halfWidth},
	}

	for _, d := range quadrantDimensions {
		q := buildQuadrant(arena, d[0], d[1], d[2], d[3])
		printMap(q, d[0], d[1], d[2], d[3])
		fmt.Println(d)
	}

	return 0
}

func main() {
	fmt.Println("Part1:", findShortestPathToCollectAllKeys())
	fmt.Println("Part2:", findShortestPathToCollectAllKeysWithFourRobots())
}
