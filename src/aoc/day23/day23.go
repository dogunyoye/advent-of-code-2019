package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type intcodeComputer struct {
	id           int
	program      []int64
	memPointer   int64
	output       []int64
	relativeBase int64
	queue        []int64
	initialised  bool
}

type nat struct {
	currentPacket  []int64
	previousPacket []int64
}

func (computer *intcodeComputer) runOpcodeForParameterMode(opcode int64, opcodeIndex int64) int {
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

		var input = int64(-1)
		if !computer.initialised {
			input = int64(computer.id)
			computer.initialised = true
		} else if len(computer.queue) != 0 {
			input = computer.queue[0]
			computer.queue = computer.queue[1:]
		}

		if paramMode0 == 0 { // position mode
			computer.program[computer.program[opcodeIndex]] = input
		} else if paramMode0 == 1 { // immediate mode
			computer.program[opcodeIndex] = input
		} else if paramMode0 == 2 { // relative mode
			computer.program[computer.program[opcodeIndex]+computer.relativeBase] = input
		}
		return 2
	}

	if paramMode0 == 0 { // position mode
		firstOperand = computer.program[computer.program[opcodeIndex]]
	} else if paramMode0 == 1 { // immediate mode
		firstOperand = computer.program[opcodeIndex]
	} else if paramMode0 == 2 { // relative mode
		firstOperand = computer.program[computer.program[opcodeIndex]+computer.relativeBase]
	} else {
		panic("Unknown first param mode: " + string(paramMode0))
	}

	// output instruction
	if opcode == 4 {
		computer.output = append(computer.output, firstOperand)
		return 2
	}

	// relative base adjustment instruction
	if opcode == 9 {
		computer.relativeBase += firstOperand
		return 2
	}

	opcodeIndex++

	if paramMode1 == 0 { // position mode
		secondOperand = computer.program[computer.program[opcodeIndex]]
	} else if paramMode1 == 1 { // immediate mode
		secondOperand = computer.program[opcodeIndex]
	} else if paramMode1 == 2 { // relative mode
		secondOperand = computer.program[computer.program[opcodeIndex]+computer.relativeBase]
	} else {
		panic("Unknown second param mode: " + string(paramMode1))
	}

	opcodeIndex++

	switch opcode {
	case 1:
		result = firstOperand + secondOperand
	case 2:
		result = firstOperand * secondOperand
	case 5:
		if firstOperand != 0 {
			computer.memPointer = secondOperand
			return 0
		}

		return 3
	case 6:
		if firstOperand == 0 {
			computer.memPointer = secondOperand
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
		computer.program[computer.program[opcodeIndex]] = result
	} else if paramMode2 == 2 {
		computer.program[computer.program[opcodeIndex]+computer.relativeBase] = result
	}

	return 4
}

func (computer *intcodeComputer) runOpcode(opcode int64, opcodeIndex int64, program []int64) int {
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
			computer.memPointer = secondOperand
			return 0
		}

		return 3
	case 6:
		if firstOperand == 0 {
			computer.memPointer = secondOperand
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

func (computer *intcodeComputer) runDiagnosticProgram(computers []*intcodeComputer, n *nat) (int64, bool) {
	opcode := computer.program[computer.memPointer]
	opcodeJump := 0

	var emptyCount = 0
	var sentPackets = 0

	for {
		opcodeJump = 2

		switch opcode {
		case 1:
			fallthrough
		case 2:
			opcodeJump = computer.runOpcode(opcode, computer.memPointer, computer.program)
		case 3:
			if !computer.initialised {
				computer.program[computer.program[computer.memPointer+1]] = int64(computer.id)
				computer.initialised = true
				break
			}

			if len(computer.queue) == 0 {
				computer.program[computer.program[computer.memPointer+1]] = -1

				// bit of a hack?
				// noticed that after the second -1 input
				// the computer has sent all packets it can
				emptyCount += 1
				if n == nil && emptyCount == 2 {
					return -1, false
				}

				// if the computer sees 2 consectutive receive request
				// where the computer has nothing to input, we can say
				// it's network is idle
				if n != nil && emptyCount == 2 {
					return -1, sentPackets == 0
				}
			} else {
				computer.program[computer.program[computer.memPointer+1]] = computer.queue[0]
				computer.queue = computer.queue[1:]
			}
		case 4:
			computer.output = append(computer.output, computer.program[computer.program[computer.memPointer+1]])
		case 5:
			fallthrough
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			opcodeJump = computer.runOpcode(opcode, computer.memPointer, computer.program)
		case 9:
			computer.relativeBase += computer.program[computer.program[computer.memPointer+1]]
		case 99:
			panic("Unexpected halt")
		default:
			opcodeJump = computer.runOpcodeForParameterMode(opcode, computer.memPointer)
		}

		computer.memPointer += int64(opcodeJump)
		opcode = computer.program[computer.memPointer]

		if len(computer.output) == 3 {
			sentPackets += 1
			address, x, y := computer.output[0], computer.output[1], computer.output[2]

			if address == 255 {
				if n == nil {
					return y, true
				}

				n.currentPacket[0] = x
				n.currentPacket[1] = y
			} else {
				computers[address].queue = append(computers[address].queue, x, y)
			}

			computer.output = make([]int64, 0)
		}
	}
}

func copyArray(array []int64) []int64 {
	arrNew := make([]int64, 0)
	arrNew = append(arrNew, array...)

	// extend program's memory
	space := make([]int64, len(array))
	arrNew = append(arrNew, space...)
	return arrNew
}

func generateProgram() []int64 {
	file, err := os.Open("../../data/day23.txt")

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

func findYValueOfFirstPacketSentToAddress255() int64 {
	computers := make([]*intcodeComputer, 50)
	for i := 0; i < 50; i++ {
		var c = intcodeComputer{i, copyArray(generateProgram()), int64(0), make([]int64, 0), int64(0), make([]int64, 0), false}
		computers[i] = &c
	}

	for {
		for i := range computers {
			val, terminate := computers[i].runDiagnosticProgram(computers, nil)
			if terminate {
				return val
			}
		}
	}
}

func findFirstYValueSentByNatTwiceInARow() int64 {
	computers := make([]*intcodeComputer, 50)
	for i := 0; i < 50; i++ {
		var c = intcodeComputer{i, copyArray(generateProgram()), int64(0), make([]int64, 0), int64(0), make([]int64, 0), false}
		computers[i] = &c
	}

	var n = nat{make([]int64, 2), make([]int64, 2)}

	for {
		idleCount := 0
		emptyQueueCount := 0

		for i := range computers {
			_, isIdle := computers[i].runDiagnosticProgram(computers, &n)
			if isIdle {
				idleCount += 1
			}
		}

		for i := range computers {
			if len(computers[i].queue) == 0 {
				emptyQueueCount += 1
			}
		}

		if idleCount == len(computers) && emptyQueueCount == len(computers) {
			computers[0].queue = append(computers[0].queue, n.currentPacket[0], n.currentPacket[1])
			if len(n.currentPacket) == len(n.previousPacket) && n.currentPacket[1] == n.previousPacket[1] {
				return n.currentPacket[1]
			}

			n.previousPacket[0] = n.currentPacket[0]
			n.previousPacket[1] = n.currentPacket[1]
			n.currentPacket = make([]int64, 2)
		}
	}
}

func main() {
	fmt.Println("Part1:", findYValueOfFirstPacketSentToAddress255())
	fmt.Println("Part2:", findFirstYValueSentByNatTwiceInARow())
}
