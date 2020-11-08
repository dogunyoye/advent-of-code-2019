package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func calculatePattern(num int, pattern []string) []int {
	var p = []int{}
	if num == 0 {
		for i := 0; i < len(pattern); i++ {
			index := (i + 1) % len(pattern)
			val, _ := strconv.Atoi(pattern[index])
			p = append(p, val)
		}
	} else {
		for i := 0; i < num; i++ {
			val, _ := strconv.Atoi("0")
			p = append(p, val)
		}

		for j := 1; j < len(pattern); j++ {

			for k := 0; k <= num; k++ {
				val, _ := strconv.Atoi(pattern[j])
				p = append(p, val)

				if len(p) == len(pattern) {
					return p
				}
			}
		}
	}

	return p
}

func runFFT(input string, pattern []string) string {
	phases := 0
	var str = ""

	for phases < 100 {
		str = ""
		for i := 0; i < len(input); i++ {
			m := calculatePattern(i, pattern)
			total := 0
			for j := 0; j < len(pattern); j++ {
				inputNum, _ := strconv.Atoi(string(input[j]))
				res := m[j] * inputNum
				total += res
			}

			totalString := strconv.Itoa(total)

			str += string(totalString[len(totalString)-1])
		}

		input = str
		phases++
	}

	return str
}

func part2(input string) []int {
	var repeated = strings.Repeat(input, 10000)
	offset, _ := strconv.Atoi(input[0:7])
	output := make([]int, 0)

	for i := offset; i < len(repeated); i++ {
		val, _ := strconv.Atoi(string(repeated[i]))
		output = append(output, val)
	}

	for phases := 0; phases < 100; phases++ {
		sum := 0
		for i := len(output) - 1; i >= 0; i-- {
			sum += output[i]
			output[i] = sum % 10
		}
	}

	return output
}

func main() {
	file, err := os.Open("../../data/day16.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var input = ""

	var pattern = []string{}
	var basePattern = []string{"0", "1", "0", "-1"}

	for scanner.Scan() {
		input = scanner.Text()
	}

	file.Close()

	for i := 0; i < len(input); i++ {
		index := i % len(basePattern)
		pattern = append(pattern, basePattern[index])
	}

	var result = runFFT(input, pattern)
	fmt.Println("Part1:", result[0:8])

	var result2 = part2(input)
	fmt.Println("Part2:", result2[0:8])
}
