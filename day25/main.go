package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"gonum.org/v1/gonum/stat/combin"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

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

// runMachine drives the adventure until the program halts (true) or
// starves for scripted input (false). Outputs stream to stdout as before.
func runMachine(m *intcode.Machine) bool {
	for {
		switch ev := m.Step(); ev.Kind {
		case intcode.Halted:
			return true
		case intcode.Output:
			fmt.Printf("%c", ev.Value)
		case intcode.NeedInput:
			if collectAllItemsIdx < len(collectAllItems) && inputIdx == 0 {
				input = collectAllItems[collectAllItemsIdx]
				collectAllItemsIdx += 1
			} else if inputIdx == 0 {
				// input, _ = reader.ReadString('\n')
				// input = input[:len(input)-2] + "\n"
				return false
			}
			m.Feed(setInput())
		case intcode.Continue:
		}
	}
}

func findPassword() {
	program, err := intcode.LoadProgram(puzzle.Path("day25.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	combos := generateItemCombos()
	for _, c := range combos {
		collectAllItemsIdx = 0

		command := generateItemCommand(c)
		collectAllItems = append(collectAllItems, command...)
		halted := runMachine(intcode.New(program))
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
