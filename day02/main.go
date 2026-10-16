package main

import (
	"fmt"
	"log"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

func main() {
	program, err := intcode.LoadProgram(puzzle.Path("day02.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	m := intcode.New(program)
	m.SetMem(1, 12)
	m.SetMem(2, 2)
	m.Run()

	fmt.Println("Part1:", m.Mem(0))

	for noun := int64(0); noun <= 99; noun++ {
		for verb := int64(0); verb <= 99; verb++ {
			m := intcode.New(program)
			m.SetMem(1, noun)
			m.SetMem(2, verb)
			m.Run()

			if m.Mem(0) == 19690720 {
				fmt.Println("Part2:", 100*m.Mem(1)+m.Mem(2))
				return
			}
		}
	}
}
