package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type moon struct {
	X    int
	Y    int
	Z    int
	xVel int
	yVel int
	zVel int
}

var initialX, initialY, initialZ = "", "", ""

func copyArray(array []moon) []moon {
	arrNew := make([]moon, 0)
	arrNew = append(arrNew, array...)
	return arrNew
}

func applyGravityAndVelocity(moons []moon) {
	for i := 0; i < len(moons)-1; i++ {
		for j := i + 1; j < len(moons); j++ {

			if moons[i].X < moons[j].X {
				moons[i].xVel++
				moons[j].xVel--
			} else if moons[i].X > moons[j].X {
				moons[i].xVel--
				moons[j].xVel++
			}

			if moons[i].Y < moons[j].Y {
				moons[i].yVel++
				moons[j].yVel--
			} else if moons[i].Y > moons[j].Y {
				moons[i].yVel--
				moons[j].yVel++
			}

			if moons[i].Z < moons[j].Z {
				moons[i].zVel++
				moons[j].zVel--
			} else if moons[i].Z > moons[j].Z {
				moons[i].zVel--
				moons[j].zVel++
			}

			if i == len(moons)-2 && j == len(moons)-1 {
				moons[j].X += moons[j].xVel
				moons[j].Y += moons[j].yVel
				moons[j].Z += moons[j].zVel
			}
		}

		moons[i].X += moons[i].xVel
		moons[i].Y += moons[i].yVel
		moons[i].Z += moons[i].zVel
	}
}

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int64, integers ...int64) int64 {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}

func searchForPreviousMoonStates(moons []moon) int64 {
	steps := int64(0)

	var exists = struct{}{}

	moonXMap := make(map[string]struct{})
	moonYMap := make(map[string]struct{})
	moonZMap := make(map[string]struct{})

	foundX, foundY, foundZ := false, false, false
	xStep, yStep, zStep := int64(0), int64(0), int64(0)

	for {
		steps++
		applyGravityAndVelocity(moons)

		moonXPos := strconv.Itoa(moons[0].X) + strconv.Itoa(moons[1].X) + strconv.Itoa(moons[2].X) + strconv.Itoa(moons[3].X)
		moonYPos := strconv.Itoa(moons[0].Y) + strconv.Itoa(moons[1].Y) + strconv.Itoa(moons[2].Y) + strconv.Itoa(moons[3].Y)
		moonZPos := strconv.Itoa(moons[0].Z) + strconv.Itoa(moons[1].Z) + strconv.Itoa(moons[2].Z) + strconv.Itoa(moons[3].Z)

		_, e1 := moonXMap[moonXPos]
		if e1 && moonXPos == initialX && !foundX {
			foundX = true
			xStep = steps
		} else {
			moonXMap[moonXPos] = exists
		}

		_, e2 := moonYMap[moonYPos]
		if e2 && moonYPos == initialY && !foundY {
			foundY = true
			yStep = steps
		} else {
			moonYMap[moonYPos] = exists
		}

		_, e3 := moonZMap[moonZPos]
		if e3 && moonZPos == initialZ && !foundZ {
			foundZ = true
			zStep = steps
		} else {
			moonZMap[moonZPos] = exists
		}

		if foundX && foundY && foundZ {
			return lcm(xStep, yStep, zStep)
		}
	}
}

func main() {
	file, err := os.Open("../../data/day12.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var moons = make([]moon, 0)

	for scanner.Scan() {
		var moonPosition = strings.Replace(scanner.Text(), "<", "", -1)
		moonPosition = strings.Replace(moonPosition, ">", "", -1)

		x, y, z := 0, 0, 0

		coords := strings.Split(moonPosition, ",")
		for i := 0; i < len(coords); i++ {
			val := strings.Split(coords[i], "=")

			switch i {
			case 0:
				x, _ = strconv.Atoi(val[1])
			case 1:
				y, _ = strconv.Atoi(val[1])
			default:
				z, _ = strconv.Atoi(val[1])
			}
		}

		var m = moon{x, y, z, 0, 0, 0}
		moons = append(moons, m)
	}

	file.Close()

	var moons2 = copyArray(moons)

	steps := 0

	for steps < 1000 {
		applyGravityAndVelocity(moons)
		steps++
	}

	total := 0

	for _, moon := range moons {
		var pot = math.Abs(float64(moon.X)) + math.Abs(float64(moon.Y)) + math.Abs(float64(moon.Z))
		var kin = math.Abs(float64(moon.xVel)) + math.Abs(float64(moon.yVel)) + math.Abs(float64(moon.zVel))

		total += (int(pot) * int(kin))
	}

	fmt.Println("Part1:", total)

	initialX = strconv.Itoa(moons2[0].X) + strconv.Itoa(moons2[1].X) + strconv.Itoa(moons2[2].X) + strconv.Itoa(moons2[3].X)
	initialY = strconv.Itoa(moons2[0].Y) + strconv.Itoa(moons2[1].Y) + strconv.Itoa(moons2[2].Y) + strconv.Itoa(moons2[3].Y)
	initialZ = strconv.Itoa(moons2[0].Z) + strconv.Itoa(moons2[1].Z) + strconv.Itoa(moons2[2].Z) + strconv.Itoa(moons2[3].Z)

	fmt.Println("Part2:", searchForPreviousMoonStates(moons2))
}
