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
	arena     map[position]rune
	queue     []positionState
	visited   map[state]struct{}
	keys      string
	steps     int
	completed bool
}

func (q *quadrant) isPathAvailable() bool {
	seen := map[positionState]struct{}{}

	for {
		currentPositionState := q.queue[0]
		_, exists := seen[currentPositionState]
		if exists {
			return false
		}

		seen[currentPositionState] = struct{}{}

		currentPosition := currentPositionState.pos
		currentValue := q.arena[currentPosition]
		currentKeys := currentPositionState.keys

		if !isDoor(currentValue) || (isDoor(currentValue) && hasKey(currentKeys, f(unicode.ToLower(currentValue)))) {
			return true
		}

		q.queue = q.queue[1:]
		q.queue = append(q.queue, currentPositionState)
	}
}

func (q *quadrant) explore(doorsInQuadrants map[rune]*quadrant) bool {

	if q.completed {
		return false
	}

	_, allKeys := findKeys(q.arena)

	for len(q.queue) != 0 {

		if !q.isPathAvailable() {
			return true
		}

		currentPositionState := q.queue[0]
		q.queue = q.queue[1:]

		currentPosition := currentPositionState.pos
		currentSteps := currentPositionState.steps
		currentKeys := currentPositionState.keys

		if hasAllArenaKeys(currentKeys, allKeys) {
			q.steps = currentSteps
			q.completed = true
			return false
		}

		neighbours := [4]position{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

		for _, n := range neighbours {
			nextPosition := position{currentPosition.i + n.i, currentPosition.j + n.j}
			nextKeys := currentKeys
			nextState := state{nextPosition, nextKeys}

			value := q.arena[nextPosition]
			if value == '#' {
				continue
			}

			_, exists := q.visited[nextState]
			if !exists {
				if isDoorKey(value) {
					nextKeys = obtainAndDistributeDoorKey(q, doorsInQuadrants, nextKeys, value)
					q.queue = append(q.queue, positionState{nextKeys, nextPosition, currentSteps + 1})
					q.visited[state{nextPosition, nextKeys}] = struct{}{}
				} else {
					if isDoor(value) || value == '.' {
						q.queue = append(q.queue, positionState{nextKeys, nextPosition, currentSteps + 1})
						q.visited[state{nextPosition, nextKeys}] = struct{}{}
					}
				}
			}
		}
	}

	if q.steps == -1 {
		panic("No solution!")
	}

	q.completed = true
	return false
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
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

func buildQuadrant(arena map[position]rune, depthStart int, depthEnd, widthStart int, widthEnd int, start position) quadrant {
	quadrantArena := make(map[position]rune)

	for i := depthStart; i < depthEnd; i++ {
		for j := widthStart; j < widthEnd; j++ {
			p := position{i, j}
			quadrantArena[p] = arena[p]
		}
	}

	emptyKeys := strings.Repeat("0", 26)
	queue := make([]positionState, 0)
	visited := map[state]struct{}{}

	queue = append(queue, positionState{emptyKeys, start, 0})
	visited[state{start, strings.Repeat("0", 26)}] = struct{}{}

	return quadrant{quadrantArena, queue, visited, emptyKeys, -1, false}
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

func findDoors(arena map[position]rune) []rune {
	doors := make([]rune, 0)
	for _, v := range arena {
		if isDoor(v) {
			doors = append(doors, v)
		}
	}
	return doors
}

func hasAllArenaKeys(currentKeys string, arenaKeys string) bool {
	for i, c := range arenaKeys {
		if c == '1' && currentKeys[i] != '1' {
			return false
		}
	}
	return true
}

func f(in rune) int {
	return int(in - 'a')
}

func obtainDoorKey(keys string, pos int) string {
	return replaceAtIndex(keys, '1', pos)
}

func obtainAndDistributeDoorKey(currentQuadrant *quadrant, doorsInQuadrants map[rune]*quadrant, keys string, doorKey rune) string {
	for k, v := range doorsInQuadrants {
		if v == currentQuadrant {
			continue
		}

		if k == unicode.ToUpper(doorKey) {
			for i := range v.queue {
				v.queue[i].keys = obtainDoorKey(v.queue[i].keys, f(doorKey))
			}
		}
	}
	return obtainDoorKey(keys, f(doorKey))
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
					nextKeys = obtainDoorKey(nextKeys, f(value))
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

	quadrants := make([]*quadrant, 0)
	doorsInQuadrants := make(map[rune]*quadrant)

	halfDepth := int(math.Round(float64(depth) / 2))
	halfWidth := int(math.Round(float64(width) / 2))

	quadrantDimensions := [][]int{
		{0, halfDepth, 0, halfWidth},
		{0, halfDepth, halfWidth - 1, width},
		{halfDepth - 1, depth, halfWidth - 1, width},
		{halfDepth - 1, depth, 0, halfWidth},
	}

	for i, d := range quadrantDimensions {
		var q = buildQuadrant(arena, d[0], d[1], d[2], d[3], startPositions[i])
		quadrants = append(quadrants, &q)

		for _, d := range findDoors(q.arena) {
			doorsInQuadrants[d] = &q
		}
	}

	for {
		allComplete := true
		for i := range quadrants {
			if quadrants[i].explore(doorsInQuadrants) {
				allComplete = false
			}
		}

		if allComplete {
			var result = 0
			for _, r := range quadrants {
				result += r.steps
			}

			return result
		}
	}
}

func main() {
	fmt.Println("Part1:", findShortestPathToCollectAllKeys())
	fmt.Println("Part2:", findShortestPathToCollectAllKeysWithFourRobots())
}
