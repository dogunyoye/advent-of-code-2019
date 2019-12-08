package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var part1 = 0

var isFound = false
var traversed = []string{}

func copyArray(array []string) []string {
	arrNew := make([]string, 0)
	arrNew = append(arrNew, array...)
	return arrNew
}

func calculateOrbiters(planet string, orbiters map[string][]string) {

	orbitingPlanets, _ := orbiters[planet]
	for _, o := range orbitingPlanets {
		part1++
		calculateOrbiters(o, orbiters)
	}
}

func calculateOrbitersPart2(startPlanet string, orbiters map[string][]string) {

	traversed = append(traversed, startPlanet)

	orbitingPlanets, _ := orbiters[startPlanet]
	for _, o := range orbitingPlanets {
		calculateOrbitersPart2(o, orbiters)
	}
}

func main() {
	file, err := os.Open("../../data/day06.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var orbit = []string{}
	var reverseOrbitMap = make(map[string][]string)

	var planets = []string{}

	for scanner.Scan() {
		orbit = strings.Split(scanner.Text(), ")")
		_, ok := reverseOrbitMap[orbit[1]]
		if !ok {
			// doesn't exist
			var dependants = []string{orbit[0]}
			reverseOrbitMap[orbit[1]] = dependants
		} else {
			reverseOrbitMap[orbit[1]] = append(reverseOrbitMap[orbit[1]], orbit[0])
		}

		planets = append(planets, orbit[1])
	}

	file.Close()

	for _, planet := range planets {
		calculateOrbiters(planet, reverseOrbitMap)
	}

	fmt.Println("Part1:", part1)

	youOrbiter := reverseOrbitMap["YOU"][0]
	sanOrbiter := reverseOrbitMap["SAN"][0]

	t1 := []string{}
	t2 := []string{}

	calculateOrbitersPart2(youOrbiter, reverseOrbitMap)
	t1 = copyArray(traversed)

	traversed = nil

	calculateOrbitersPart2(sanOrbiter, reverseOrbitMap)
	t2 = copyArray(traversed)

	distanceYou := 0
	distanceSanta := 0

	for i, x := range t1 {
		for k, y := range t2 {
			if x == y {
				distanceYou = i
				distanceSanta = k
				isFound = true
			}
		}

		if isFound {
			break
		}
	}

	fmt.Println("Part2:", distanceYou+distanceSanta)
}
