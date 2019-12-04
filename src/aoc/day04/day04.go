package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("../../data/day04.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	input := []string{}

	for scanner.Scan() {
		input = strings.Split(scanner.Text(), "-")
	}

	file.Close()

	lowerLimit, _ := strconv.Atoi(input[0])
	upperLimit, _ := strconv.Atoi(input[1])

	// using a map to eliminate duplicates
	passwordCandidates := make(map[int]bool)
	result := 0

	passwordCandidates2 := []string{}

	for i := lowerLimit; i <= upperLimit; i++ {
		number := strconv.Itoa(i)
		for j := 0; j < len(number)-1; j++ {
			if number[j] == number[j+1] {
				passwordCandidates[i] = true
			}
		}
	}

	for key := range passwordCandidates {
		isCandidate := true
		candidate := strconv.Itoa(key)

		for k := 0; k < len(candidate)-1; k++ {
			// if the left number is greater than the right number
			// terminate. Rule broken
			if candidate[k] > candidate[k+1] {
				isCandidate = false
				break
			}
		}

		if isCandidate {
			passwordCandidates2 = append(passwordCandidates2, candidate)
			result++
		}
	}

	fmt.Println("Part1:", result)

	part2Result := 0

	for _, part2Candidate := range passwordCandidates2 {
		dupCounter := 0
		isPart2Candidate := false

		for k := 0; k < len(part2Candidate)-1; k++ {
			if part2Candidate[k] == part2Candidate[k+1] {
				dupCounter++
			} else {
				if dupCounter == 1 {
					break
				} else {
					dupCounter = 0
				}
			}

			if dupCounter == 1 {
				isPart2Candidate = true
			} else {
				isPart2Candidate = false
			}
		}

		if isPart2Candidate {
			part2Result++
		}

	}

	fmt.Println("Part2:", part2Result)
}
