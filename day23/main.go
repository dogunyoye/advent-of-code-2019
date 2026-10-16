package main

import (
	"fmt"
	"log"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

type intcodeComputer struct {
	id          int
	vm          *intcode.Machine
	output      []int64
	queue       []int64
	initialised bool
}

type nat struct {
	currentPacket  []int64
	previousPacket []int64
}

func (computer *intcodeComputer) runDiagnosticProgram(computers []*intcodeComputer, n *nat) (int64, bool) {
	var emptyCount = 0
	var sentPackets = 0

	for {
		switch ev := computer.vm.Step(); ev.Kind {
		case intcode.Halted:
			panic("Unexpected halt")
		case intcode.Output:
			computer.output = append(computer.output, ev.Value)
		case intcode.NeedInput:
			computer.feedInput(&emptyCount)
			if emptyCount == 2 {
				// bit of a hack?
				// noticed that after the second -1 input
				// the computer has sent all packets it can
				if n == nil {
					return -1, false
				}

				// if the computer sees 2 consectutive receive request
				// where the computer has nothing to input, we can say
				// it's network is idle
				return -1, sentPackets == 0
			}
		case intcode.Continue:
		}

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

// feedInput supplies the pending input instruction, mirroring the
// original split behaviour: only a bare position-mode opcode 3 (word
// exactly 3) advances the idle counter; mode-prefixed inputs never do.
// emptyCount is incremented only when a bare input starves.
func (computer *intcodeComputer) feedInput(emptyCount *int) {
	bare := computer.vm.Mem(computer.vm.IP()) == 3

	var input = int64(-1)
	if !computer.initialised {
		input = int64(computer.id)
		computer.initialised = true
	} else if len(computer.queue) != 0 {
		input = computer.queue[0]
		computer.queue = computer.queue[1:]
	} else if bare {
		*emptyCount += 1
	}
	computer.vm.Feed(input)
}

func findYValueOfFirstPacketSentToAddress255() int64 {
	program, err := intcode.LoadProgram(puzzle.Path("day23.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	computers := make([]*intcodeComputer, 50)
	for i := 0; i < 50; i++ {
		computers[i] = &intcodeComputer{
			id:          i,
			vm:          intcode.New(program),
			output:      make([]int64, 0),
			queue:       make([]int64, 0),
			initialised: false,
		}
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
	program, err := intcode.LoadProgram(puzzle.Path("day23.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	computers := make([]*intcodeComputer, 50)
	for i := 0; i < 50; i++ {
		computers[i] = &intcodeComputer{
			id:          i,
			vm:          intcode.New(program),
			output:      make([]int64, 0),
			queue:       make([]int64, 0),
			initialised: false,
		}
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
