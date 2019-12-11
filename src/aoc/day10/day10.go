package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

type position struct {
	X int
	Y int
}

func scanAsteroids(asteroids []position) (int, position) {
	mostInView := 0
	mostInViewPosition := position{}

	for _, currentAsteroid := range asteroids {
		var asteroidInView = make(map[float64]position)

		for _, otherAsteroid := range asteroids {

			if (currentAsteroid.X == otherAsteroid.X) && (currentAsteroid.Y == otherAsteroid.Y) {
				// skip evaluating against the asteroid I'm on obviously
				continue
			}

			xDiff := (otherAsteroid.X - currentAsteroid.X)
			yDiff := (otherAsteroid.Y - currentAsteroid.Y)

			angle := math.Atan2(float64(-yDiff), float64(xDiff))

			_, ok := asteroidInView[angle]
			if !ok {
				asteroidInView[angle] = otherAsteroid
			} else {
				// calculate the distance to this position
				// if it's closer, replace the current position
				// in the map with this.
				pos := asteroidInView[angle]
				posXDiff := (pos.X - currentAsteroid.X)
				posYDiff := (pos.Y - currentAsteroid.Y)

				m1 := math.Abs(float64(xDiff)) + math.Abs(float64(yDiff))
				m2 := math.Abs(float64(posXDiff)) + math.Abs(float64(posYDiff))

				if m2 > m1 {
					delete(asteroidInView, angle)
					asteroidInView[angle] = otherAsteroid
				}
			}
		}

		inView := len(asteroidInView)

		if inView > mostInView {
			mostInView = inView
			mostInViewPosition = currentAsteroid
		}
	}

	return mostInView, mostInViewPosition
}

func scanAsteroidsAndDestroy(station position, asteroids []position) position {
	var asteroidInView = make(map[float64]position)

	for _, otherAsteroid := range asteroids {

		if (station.X == otherAsteroid.X) && (station.Y == otherAsteroid.Y) {
			// skip evaluating against the asteroid I'm on obviously
			continue
		}

		xDiff := (otherAsteroid.X - station.X)
		yDiff := (otherAsteroid.Y - station.Y)

		angle := math.Atan2(float64(-yDiff), float64(xDiff))

		_, ok := asteroidInView[angle]
		if !ok {
			asteroidInView[angle] = otherAsteroid
		} else {
			// calculate the distance to this position
			// if it's closer, replace the current position
			// in the map with this.
			pos := asteroidInView[angle]
			posXDiff := (pos.X - station.X)
			posYDiff := (pos.Y - station.Y)

			m1 := math.Abs(float64(xDiff)) + math.Abs(float64(yDiff))
			m2 := math.Abs(float64(posXDiff)) + math.Abs(float64(posYDiff))

			if m2 > m1 {
				delete(asteroidInView, angle)
				asteroidInView[angle] = otherAsteroid
			}
		}
	}

	angles := []float64{}
	anglesInUpOrder := []float64{}

	for k := range asteroidInView {
		angles = append(angles, k)
	}

	sort.Float64s(angles)

	for i, j := 0, len(angles)-1; i < j; i, j = i+1, j-1 {
		angles[i], angles[j] = angles[j], angles[i]
	}

	// create a new list and append angles starting in the up
	// direction, clockwise
	for i, a := range angles {
		if a <= math.Atan2(float64(1), float64(0)) && a > math.Atan2(float64(-1), float64(0)) {
			anglesInUpOrder = append(anglesInUpOrder, angles[i:]...)
			anglesInUpOrder = append(anglesInUpOrder, angles[:i]...)
			break
		}
	}

	for i, p := range anglesInUpOrder {
		if i == 199 {
			return asteroidInView[p]
		}
	}

	return station
}

func main() {
	file, err := os.Open("../../data/day10.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	row := 0
	index := 0
	asteroids := []position{}

	for scanner.Scan() {
		line := scanner.Text()
		for i := 0; i < len(line); i++ {
			if line[i] == '#' {
				asteroids = append(asteroids, position{i, row})
				index++
			}
		}

		row++
	}

	file.Close()

	mostInView, mostInViewPosition := scanAsteroids(asteroids)

	fmt.Println("Part1:", mostInView, "( at position:", mostInViewPosition, ")")

	a := scanAsteroidsAndDestroy(mostInViewPosition, asteroids)

	fmt.Println("Part2:", ((a.X * 100) + a.Y))
}
