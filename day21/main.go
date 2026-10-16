package main

import (
	"fmt"
	"log"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

var positionIndex = 0

// https://todd.ginsberg.com/post/advent-of-code/2019/day21/
var walk = "NOT A J\nNOT B T\nAND D T\nOR T J\nNOT C T\nOR T J\nAND D J\nWALK\n"
var run = "NOT C J\nAND D J\nNOT H T\nNOT T T\nOR E T\nAND T J\nNOT A T\nOR T J\nNOT B T\nNOT T T\nOR E T\nNOT T T\nOR T J\nRUN\n"

func setInput(part2 bool) int64 {
	var input = int64(0)
	if !part2 {
		input = int64(walk[positionIndex])
	} else {
		input = int64(run[positionIndex])
	}

	positionIndex++
	return input
}

func main() {
	program, err := intcode.LoadProgram(puzzle.Path("day21.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	run := func(part2 bool) int64 {
		positionIndex = 0
		var last int64
		m := intcode.New(program)
		m.In = func() int64 { return setInput(part2) }
		m.Out = func(v int64) { last = v }
		m.Run()
		return last
	}

	fmt.Println("Part1:", run(false))
	fmt.Println("Part2:", run(true))
}
