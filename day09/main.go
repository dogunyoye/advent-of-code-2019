package main

import (
	"fmt"
	"log"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

func main() {
	program, err := intcode.LoadProgram(puzzle.Path("day09.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	run := func(input int64) {
		m := intcode.New(program)
		m.In = func() int64 { return input }
		m.Out = func(v int64) { fmt.Println("output:", v) }
		m.Run()
	}

	run(1)
	run(2)
}
