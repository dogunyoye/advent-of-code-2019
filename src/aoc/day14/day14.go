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

// 10 ORE => 10 A
type oreReaction struct {
	numberOfOutput int
	numberOfOre    int
}

type ingredient struct {
	number int
	letter string
}

// 7 A, 1 D => 1 E
type formulaReaction struct {
	numberOfOutput int
	ingredients    []ingredient
}

var oreFormulas = make(map[string]oreReaction)
var normalFormulas = make(map[string]formulaReaction)
var fuelFormula = ""

var tempOre = 0
var limitReached = false

var wasted = make(map[string]int)

func calculateIngredientsNeeded(numNeeded int, letter string, total *int) {

	var output = 0

	_, isOreFormula := oreFormulas[letter]
	if isOreFormula {
		var f = oreFormulas[letter]
		output = f.numberOfOutput

	} else {
		var f = normalFormulas[letter]
		output = f.numberOfOutput
	}

	var iterations = 0

	_, wasteExists := wasted[letter]
	if wasteExists {
		if wasted[letter] <= numNeeded {
			numNeeded -= wasted[letter]
			wasted[letter] = 0
		} else {
			wasted[letter] -= numNeeded
			numNeeded = 0
		}
	} else {
		wasted[letter] = 0
	}

	if numNeeded == 0 {
		return
	}

	if numNeeded <= output {
		iterations = 1
	} else {
		iterations = int(math.Ceil(float64(numNeeded) / float64(output)))
	}

	// check if waste generated
	excess := (iterations * output) - numNeeded
	wasted[letter] += excess

	if isOreFormula {

		var needed = numNeeded
		var current = 0
		var r = oreFormulas[letter]
		var ore = 0

		for current < needed {
			current += r.numberOfOutput
			ore += r.numberOfOre
		}

		*total += ore
	} else {
		var f = normalFormulas[letter]

		for _, ingredient := range f.ingredients {
			calculateIngredientsNeeded(ingredient.number*iterations, ingredient.letter, total)
		}
	}

}

// Brute force solution - must be adapted into something smarter
func calculateMaxFuelFromOneBillionOre(numNeeded int, letter string, total *int) {
	var output = 0

	_, isOreFormula := oreFormulas[letter]
	if isOreFormula {
		var f = oreFormulas[letter]
		output = f.numberOfOutput

	} else {
		var f = normalFormulas[letter]
		output = f.numberOfOutput
	}

	var iterations = 0

	_, wasteExists := wasted[letter]
	if wasteExists {
		if wasted[letter] <= numNeeded {
			numNeeded -= wasted[letter]
			wasted[letter] = 0
		} else {
			wasted[letter] -= numNeeded
			numNeeded = 0
		}
	} else {
		wasted[letter] = 0
	}

	if numNeeded == 0 {
		return
	}

	if numNeeded <= output {
		iterations = 1
	} else {
		iterations = int(math.Ceil(float64(numNeeded) / float64(output)))
	}

	// check if waste generated
	excess := (iterations * output) - numNeeded
	wasted[letter] += excess

	if isOreFormula {

		var needed = numNeeded
		var current = 0
		var r = oreFormulas[letter]

		var ore = 0

		for current < needed {
			current += r.numberOfOutput

			tempOre += r.numberOfOre
			if tempOre < 1000000000000 {
				ore += r.numberOfOre
			} else {
				limitReached = true
				return
			}
		}

		*total += ore
	} else {
		var f = normalFormulas[letter]

		for _, ingredient := range f.ingredients {
			calculateMaxFuelFromOneBillionOre(ingredient.number*iterations, ingredient.letter, total)
		}
	}
}

func main() {
	file, err := os.Open("../../data/day14.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "FUEL") {
			fuelFormula = line
		} else {

			var input = strings.Split(line, " => ")

			if strings.Contains(line, "ORE") {
				leftSide := strings.Split(input[0], " ")
				rightSide := strings.Split(input[1], " ")

				oreNumber, _ := strconv.Atoi(leftSide[0])
				outputNumber, _ := strconv.Atoi(rightSide[0])

				oreFormulas[rightSide[1]] = oreReaction{outputNumber, oreNumber}
			} else {
				leftSide := strings.Split(input[0], ", ")
				ingredients := []ingredient{}
				for _, chems := range leftSide {
					c := strings.Split(chems, " ")
					num, _ := strconv.Atoi(c[0])
					letter := c[1]

					var i = ingredient{num, letter}
					ingredients = append(ingredients, i)
				}

				rightSide := strings.Split(input[1], " ")
				outputNumber, _ := strconv.Atoi(rightSide[0])
				outputLetter := rightSide[1]

				var f = formulaReaction{outputNumber, ingredients}
				normalFormulas[outputLetter] = f
			}
		}

	}

	file.Close()

	total := 0

	left := strings.Split(strings.Split(fuelFormula, " => ")[0], ", ")
	for _, chems := range left {
		c := strings.Split(chems, " ")
		num, _ := strconv.Atoi(c[0])
		letter := c[1]

		calculateIngredientsNeeded(num, letter, &total)
	}

	fmt.Println("Part1:", total)
	total = 0
	fuelProduced := 0

	for {
		for _, chems := range left {
			c := strings.Split(chems, " ")
			num, _ := strconv.Atoi(c[0])
			letter := c[1]

			calculateMaxFuelFromOneBillionOre(num, letter, &total)
		}

		if limitReached {
			break
		}

		fuelProduced++
	}

	fmt.Println("Part2:", fuelProduced-1)

}
