package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/stat/combin"
)

var memPointer = int64(0)
var output = int64(0)
var relativeBase = int64(0)

var inputIdx = 0
var input string

// use this to draw out a map of the points of interest:
// all collectable items and the security checkpoint
//
//lint:ignore U1000 Ignore unused function for debugging/playing the game!
var reader = bufio.NewReader(os.Stdin)

var collectAllItemsIdx = 0

// HARDCODED
// used pen, paper and manual inputs to draw a map of the items,
// their locations and the location of security checkpoint.
// These instruction collects all carryable items and takes you
// to the security checkpoint
var collectAllItems = []string{
	"east\n",
	"take manifold\n",
	"south\n",
	"take whirled peas\n",
	"north\n",
	"west\n",
	"south\n",
	"take space heater\n",
	"south\n",
	"take dark matter\n",
	"north\n",
	"east\n",
	"north\n",
	"west\n",
	"south\n",
	"take antenna\n",
	"north\n",
	"east\n",
	"south\n",
	"east\n",
	"take bowl of rice\n",
	"north\n",
	"take klein bottle\n",
	"north\n",
	"take spool of cat6\n",
	"west\n",
	"drop manifold\n",
	"drop whirled peas\n",
	"drop space heater\n",
	"drop dark matter\n",
	"drop antenna\n",
	"drop bowl of rice\n",
	"drop klein bottle\n",
	"drop spool of cat6\n",
}

// all collectable items in the map
var items = []string{
	"manifold",
	"whirled peas",
	"space heater",
	"dark matter",
	"antenna",
	"bowl of rice",
	"klein bottle",
	"spool of cat6",
}

func generateItemCombos() [][]string {
	n := len(items)
	k := 4

	combinations := make([][]string, 0)
	list := combin.Combinations(n, k)

	for _, i := range list {
		combo := make([]string, 0)
		for _, j := range i {
			combo = append(combo, items[j])
		}
		combinations = append(combinations, combo)
	}

	return combinations
}

func generateItemCommand(items []string) []string {
	commands := make([]string, 0)
	for _, i := range items {
		commands = append(commands, "take "+i+"\n")
	}
	commands = append(commands, "north\n")
	return commands
}

func setInput() int64 {
	res := int64(input[inputIdx])
	inputIdx += 1
	if inputIdx >= len(input) {
		input = ""
		inputIdx = 0
	}
	return res
}

func copyArray(array []int64) []int64 {
	arrNew := make([]int64, 0)
	arrNew = append(arrNew, array...)

	// extend program's memory
	space := make([]int64, len(array)*1000)
	arrNew = append(arrNew, space...)
	return arrNew
}

func generateProgram() []int64 {
	file, err := os.Open("../../data/day25.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var input = []string{}

	for scanner.Scan() {
		input = strings.Split(scanner.Text(), ",")
	}

	var program = []int64{}

	for _, i := range input {
		j, err := strconv.Atoi(i)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		program = append(program, int64(j))
	}

	file.Close()

	return program
}

func runOpcodeForParameterMode(opcode int64, opcodeIndex int64, program []int64) (int, bool) {
	opcodeAsString := strconv.Itoa(int(opcode))

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

	o, _ := strconv.Atoi(code)
	opcode = int64(o)

	paramMode0, _ := strconv.Atoi(p1)
	paramMode1, _ := strconv.Atoi(p2)
	paramMode2, _ := strconv.Atoi(p3)

	var firstOperand = int64(0)
	var secondOperand = int64(0)

	var result = int64(0)

	opcodeIndex++

	// input instruction
	if opcode == 3 {
		if collectAllItemsIdx < len(collectAllItems) && inputIdx == 0 {
			input = collectAllItems[collectAllItemsIdx]
			collectAllItemsIdx += 1
		} else if inputIdx == 0 {
			// input, _ = reader.ReadString('\n')
			// input = input[:len(input)-2] + "\n"
			return -1, false
		}
		if paramMode0 == 0 { // position mode
			program[program[opcodeIndex]] = setInput()
		} else if paramMode0 == 1 { // immediate mode
			program[opcodeIndex] = setInput()
		} else if paramMode0 == 2 { // relative mode
			program[program[opcodeIndex]+relativeBase] = setInput()
		}
		return 2, true
	}

	if paramMode0 == 0 { // position mode
		firstOperand = program[program[opcodeIndex]]
	} else if paramMode0 == 1 { // immediate mode
		firstOperand = program[opcodeIndex]
	} else if paramMode0 == 2 { // relative mode
		firstOperand = program[program[opcodeIndex]+relativeBase]
	} else {
		fmt.Println("Unknown first param mode:", paramMode0)
	}

	// output instruction
	if opcode == 4 {
		output = firstOperand
		fmt.Printf("%c", output)
		return 2, true
	}

	// relative base adjustment instruction
	if opcode == 9 {
		relativeBase += firstOperand
		return 2, true
	}

	opcodeIndex++

	if paramMode1 == 0 { // position mode
		secondOperand = program[program[opcodeIndex]]
	} else if paramMode1 == 1 { // immediate mode
		secondOperand = program[opcodeIndex]
	} else if paramMode1 == 2 { // relative mode
		secondOperand = program[program[opcodeIndex]+relativeBase]
	} else {
		fmt.Println("Unknown second param mode:", paramMode1)
	}

	opcodeIndex++

	switch opcode {
	case 1:
		result = firstOperand + secondOperand
	case 2:
		result = firstOperand * secondOperand
	case 5:
		if firstOperand != 0 {
			memPointer = secondOperand
			return 0, true
		}

		return 3, true
	case 6:
		if firstOperand == 0 {
			memPointer = secondOperand
			return 0, true
		}

		return 3, true
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
	} else if paramMode2 == 2 {
		program[program[opcodeIndex]+relativeBase] = result
	}

	return 4, true
}

func runOpcode(opcode int64, opcodeIndex int64, program []int64) int {
	var firstOperand = program[program[opcodeIndex+1]]
	var secondOperand = program[program[opcodeIndex+2]]
	var resultIndex = program[opcodeIndex+3]

	var result = int64(0)

	switch opcode {
	case 1:
		result = firstOperand + secondOperand
	case 2:
		result = firstOperand * secondOperand
	case 5:
		if firstOperand != 0 {
			memPointer = secondOperand
			return 0
		}

		return 3
	case 6:
		if firstOperand == 0 {
			memPointer = secondOperand
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

func runDiagnosticProgram(program []int64) bool {
	opcode := program[memPointer]

	opcodeJump := 0
	continueProgram := true

	for {
		opcodeJump = 2

		switch opcode {
		case 1:
			fallthrough
		case 2:
			opcodeJump = runOpcode(opcode, memPointer, program)
		case 3:
			if collectAllItemsIdx < len(collectAllItems) && inputIdx == 0 {
				input = collectAllItems[collectAllItemsIdx]
				collectAllItemsIdx += 1
			} else if inputIdx == 0 {
				// input, _ = reader.ReadString('\n')
				// input = input[:len(input)-2] + "\n"
				return false
			}
			program[program[memPointer+1]] = setInput()
		case 4:
			output = program[program[memPointer+1]]
			fmt.Printf("%c", output)
		case 5:
			fallthrough
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			opcodeJump = runOpcode(opcode, memPointer, program)
		case 9:
			relativeBase += program[program[memPointer+1]]
		case 99:
			return true
		default:
			opcodeJump, continueProgram = runOpcodeForParameterMode(opcode, memPointer, program)
			if !continueProgram {
				return false
			}
		}

		memPointer += int64(opcodeJump)
		opcode = program[memPointer]
	}
}

func findPassword() {
	combos := generateItemCombos()
	for _, c := range combos {
		memPointer = int64(0)
		output = int64(0)
		relativeBase = int64(0)
		collectAllItemsIdx = 0

		command := generateItemCommand(c)
		collectAllItems = append(collectAllItems, command...)
		halted := runDiagnosticProgram(copyArray(generateProgram()))
		if halted {
			return
		}

		// remove the last 5 instructions, ie the 4 items we just
		// picked up and the movement to the security checkpoint
		collectAllItems = collectAllItems[:len(collectAllItems)-5]
	}

	panic("no solution found!")
}

func main() {
	findPassword()
}
