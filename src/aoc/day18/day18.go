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

func printMap(arena map[position]rune, depth int, width int) {
	for i := 0; i < depth; i++ {
		line := ""
		for j := 0; j < width; j++ {
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

func main() {
	fmt.Println("Part1:", findShortestPathToCollectAllKeys())
}
