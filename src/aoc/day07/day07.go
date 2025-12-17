package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var output = 0

var amplifierOutput = 0

type amplifier struct {
	code   string
	input  chan int
	output chan int
}

func copyArray(array []int) []int {
	arrNew := make([]int, 0)
	arrNew = append(arrNew, array...)
	return arrNew
}

func generateCombinations(inputs []string) []string {
	var combinations = []string{}

	for _, a := range inputs {
		for _, b := range inputs {
			for _, c := range inputs {
				for _, d := range inputs {
					for _, e := range inputs {
						duplicates := make(map[string]bool)

						duplicates[a] = true

						_, ok0 := duplicates[b]
						if !ok0 {
							duplicates[b] = true
						}
						_, ok1 := duplicates[c]
						if !ok1 {
							duplicates[c] = true
						}
						_, ok2 := duplicates[d]
						if !ok2 {
							duplicates[d] = true
						}

						_, ok3 := duplicates[e]

						if ok0 || ok1 || ok2 || ok3 {
							continue
						} else {
							combination := a + b + c + d + e
							combinations = append(combinations, combination)
						}

					}
				}
			}
		}
	}

	return combinations
}

func generateProgram() []int {
	file, err := os.Open("../../data/day07.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var input = []string{}

	for scanner.Scan() {
		input = strings.Split(scanner.Text(), ",")
	}

	var program = []int{}

	for _, i := range input {
		j, err := strconv.Atoi(i)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		program = append(program, j)
	}

	file.Close()

	return program
}

func runOpcodeForParameterMode(opcode int, memPointer *int, program []int) int {
	var opcodeIndex = *memPointer
	opcodeAsString := strconv.Itoa(opcode)

	i := 5 - len(opcodeAsString)

	padder := ""
	for x := 0; x < i; x++ {
		padder += "0"
	}

	// pad with 0
	opcodeAsString = padder + opcodeAsString

	code := string(opcodeAsString[3]) + string(opcodeAsString[4])
	p1 := string(opcodeAsString[2])
	p2 := string(opcodeAsString[1])
	p3 := string(opcodeAsString[0])

	opcode, _ = strconv.Atoi(code)
	paramMode0, _ := strconv.Atoi(p1)
	paramMode1, _ := strconv.Atoi(p2)
	paramMode2, _ := strconv.Atoi(p3)

	firstOperand := 0
	secondOperand := 0

	result := 0

	opcodeIndex++

	if paramMode0 == 0 {
		firstOperand = program[program[opcodeIndex]]
	} else {
		firstOperand = program[opcodeIndex]
	}

	if opcode == 4 {
		if paramMode0 == 0 {
			output = program[firstOperand]
		} else {
			output = firstOperand
		}

		amplifierOutput = output
		return 2
	}

	opcodeIndex++

	if paramMode1 == 0 {
		secondOperand = program[program[opcodeIndex]]
	} else {
		secondOperand = program[opcodeIndex]
	}

	opcodeIndex++

	switch opcode {
	case 1:
		result = firstOperand + secondOperand
	case 2:
		result = firstOperand * secondOperand
	case 5:
		if firstOperand != 0 {
			*memPointer = secondOperand
			return 0
		}

		return 3
	case 6:
		if firstOperand == 0 {
			*memPointer = secondOperand
			return 0
		}

		return 3
	case 7:
		if firstOperand < secondOperand {
			result = 1
		} else {
			result = 0
		}
	case 8:
		if firstOperand == secondOperand {
			result = 1
		} else {
			result = 0
		}
	}

	if paramMode2 == 0 {
		program[program[opcodeIndex]] = result
	}

	return 4
}

func runOpcode(opcode int, memPointer *int, program []int) int {
	firstOperand := program[program[*memPointer+1]]
	secondOperand := program[program[*memPointer+2]]
	resultIndex := program[*memPointer+3]

	result := 0

	switch opcode {
	case 1:
		result = firstOperand + secondOperand
	case 2:
		result = firstOperand * secondOperand
	case 5:
		if firstOperand != 0 {
			*memPointer = secondOperand
			return 0
		}

		return 3
	case 6:
		if firstOperand == 0 {
			*memPointer = secondOperand
			return 0
		}

		return 3
	case 7:
		if firstOperand < secondOperand {
			result = 1
		} else {
			result = 0
		}
	case 8:
		if firstOperand == secondOperand {
			result = 1
		} else {
			result = 0
		}
	}

	program[resultIndex] = result

	return 4
}

func runDiagnosticProgram(phaseSetting int, program []int) {
	var memPointer = 0
	var numOfInputs = 0

	opcode := program[memPointer]

	opcodeJump := 0

	for {
		opcodeJump = 2

		switch opcode {
		case 1:
			fallthrough
		case 2:
			opcodeJump = runOpcode(opcode, &memPointer, program)
		case 3:
			numOfInputs++

			// on every second input request, we take the previous amplifiers result
			if numOfInputs%2 == 0 {
				program[program[memPointer+1]] = amplifierOutput
			} else {
				program[program[memPointer+1]] = phaseSetting
			}
		case 4:
			amplifierOutput = program[program[memPointer+1]]
		case 5:
			fallthrough
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			opcodeJump = runOpcode(opcode, &memPointer, program)
		case 99:
			return
		default:
			opcodeJump = runOpcodeForParameterMode(opcode, &memPointer, program)
		}

		memPointer += opcodeJump
		opcode = program[memPointer]
	}
}

func runDiagnosticProgramPart2(phaseSetting int, amp amplifier, wg *sync.WaitGroup, completed *int32, program []int) {
	var memPointer = 0
	var numOfInputs = 0

	var ampAInitialised = false

	opcode := program[memPointer]

	opcodeJump := 0

	for {
		opcodeJump = 2

		switch opcode {
		case 1:
			fallthrough
		case 2:
			opcodeJump = runOpcode(opcode, &memPointer, program)
		case 3:
			numOfInputs++

			if numOfInputs == 1 {
				program[program[memPointer+1]] = phaseSetting
			} else {
				if amp.code == "A" && !ampAInitialised {
					program[program[memPointer+1]] = 0
					ampAInitialised = true
				} else {
					in := <-amp.input
					program[program[memPointer+1]] = in
				}
			}
		case 4:
			if amp.code == "E" {
				amplifierOutput = program[program[memPointer+1]]
				// super hacky sleep to ensure the last amplifier doesn't
				// race ahead of the previous amplifier (D) to ensure a
				// correct value for completed
				time.Sleep(1 * time.Microsecond)
			}

			if atomic.LoadInt32(completed) > 1 {
				amp.output <- program[program[memPointer+1]]
			}

		case 5:
			fallthrough
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			opcodeJump = runOpcode(opcode, &memPointer, program)
		case 99:
			atomic.AddInt32(completed, -1)
			wg.Done()

			return
		default:
			opcodeJump = runOpcodeForParameterMode(opcode, &memPointer, program)
		}

		memPointer += opcodeJump
		opcode = program[memPointer]
	}
}

func main() {

	var program1 = copyArray(generateProgram())

	combinations := generateCombinations([]string{"0", "1", "2", "3", "4"})

	max := 0
	maxCombo := ""

	for _, combo := range combinations {
		var program = copyArray(program1)
		for i := 0; i < len(combo); i++ {
			phaseSetting, _ := strconv.Atoi(string(combo[i]))
			runDiagnosticProgram(phaseSetting, program)
		}

		if amplifierOutput > max {
			max = amplifierOutput
			maxCombo = combo
		}

		amplifierOutput = 0
	}

	fmt.Println("Part1:", max, maxCombo)
	amplifierOutput = 0
	max = 0
	maxCombo = ""

	combinationsPart2 := generateCombinations([]string{"5", "6", "7", "8", "9"})

	var aChannel = make(chan int)
	var bChannel = make(chan int)
	var cChannel = make(chan int)
	var dChannel = make(chan int)
	var eChannel = make(chan int)

	var ampA = amplifier{"A", aChannel, bChannel}
	var ampB = amplifier{"B", bChannel, cChannel}
	var ampC = amplifier{"C", cChannel, dChannel}
	var ampD = amplifier{"D", dChannel, eChannel}
	var ampE = amplifier{"E", eChannel, aChannel}

	var amplifiers = [5]amplifier{ampA, ampB, ampC, ampD, ampE}

	for _, comboPart2 := range combinationsPart2 {
		var completed = int32(5)
		var wg sync.WaitGroup

		for i := 0; i < len(comboPart2); i++ {
			phaseSetting, _ := strconv.Atoi(string(comboPart2[i]))
			wg.Add(1)
			go runDiagnosticProgramPart2(phaseSetting, amplifiers[i], &wg, &completed, copyArray(program1))
		}

		// block here until last amplifer halts
		wg.Wait()

		if amplifierOutput > max {
			max = amplifierOutput
			maxCombo = comboPart2
		}

		amplifierOutput = 0
	}

	fmt.Println("Part2:", max, maxCombo)
}
