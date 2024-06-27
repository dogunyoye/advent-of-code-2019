package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"unicode"
)

type position struct {
	i int
	j int
}

type positionState struct {
	key   rune
	pos   position
	steps int
}

type state struct {
	pos  position
	keys int64
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

func unlockDoorIfPresent(key rune, arena map[position]rune) map[position]rune {
	newArena := make(map[position]rune)
	for k, v := range arena {
		p := position{k.i, k.j}
		if v == key || v == unicode.ToUpper(key) {
			newArena[p] = '.'
		} else {
			newArena[p] = v
		}
	}

	return newArena
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

func findKeys(arena map[position]rune) int {
	keys := 0
	for _, v := range arena {
		if isDoorKey(v) {
			keys += 1
		}
	}
	return keys
}

func f(in rune) int {
	return int(in - 'a' + 1)
}

func setBit(n int64, pos int) int64 {
	n |= (1 << pos)
	return n
}

func collectKeys(arena map[position]rune, depth int, width int, startState positionState, heldKeys int64, memo map[state]int) int {

	res, exists := memo[state{startState.pos, heldKeys}]
	if exists && res < startState.steps {
		if res < startState.steps {
			return res
		} else {
			memo[state{startState.pos, heldKeys}] = startState.steps
		}
		//fmt.Println("hit:", startState.pos, startState.steps, res)
	}

	if findKeys(arena) == 0 {
		return startState.steps
	}

	queue := make([]positionState, 0)
	visited := map[position]struct{}{}

	foundKeys := make([]positionState, 0)

	queue = append(queue, startState)
	visited[startState.pos] = struct{}{}

	distance := math.MaxUint32

	for len(queue) != 0 {

		currentPositionState := queue[0]
		queue = queue[1:]

		currentPosition := currentPositionState.pos
		currentSteps := currentPositionState.steps

		neighbours := [4]position{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

		for _, n := range neighbours {
			nextPosition := position{currentPosition.i + n.i, currentPosition.j + n.j}
			value := arena[nextPosition]
			_, exists := visited[nextPosition]

			if value == '#' || isDoor(arena[nextPosition]) {
				continue
			}

			if !exists {
				visited[nextPosition] = struct{}{}
				if isDoorKey(value) {
					foundKeys = append(foundKeys, positionState{value, nextPosition, currentSteps + 1})
				} else {
					queue = append(queue, positionState{-1, nextPosition, currentSteps + 1})
				}
			}
		}
	}

	for _, k := range foundKeys {
		//fmt.Printf("%c\n", k)
		//path += string(k)
		newArena := unlockDoorIfPresent(k.key, arena)
		var newHeldKeys = setBit(heldKeys, f(k.key))
		newState := state{k.pos, newHeldKeys}
		//printMap(newArena, depth, width)

		result := collectKeys(newArena, depth, width, k, newHeldKeys, memo)
		//path = path[:len(path)-1]
		memo[newState] = result

		//fmt.Println("RESULT:", result)
		if result < distance {
			distance = result
		}
	}

	return distance
}

func findShortestPathToCollectAllKeys() int {
	arena, depth, width, start := buildMap()
	printMap(arena, depth, width)

	memo := make(map[state]int)
	keys := int64(0)

	result := collectKeys(arena, depth, width, positionState{-1, start, 0}, keys, memo)
	fmt.Println("debug:", memo)
	return result
}

func main() {
	fmt.Println("Part1:", findShortestPathToCollectAllKeys())
}
