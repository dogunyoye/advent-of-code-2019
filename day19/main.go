package main

import (
	"fmt"
	"log"
	"time"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

// isPulled runs the tractor-beam program for coordinate (x, y).
func isPulled(program []int64, x, y int) bool {
	m := intcode.New(program)
	in := [2]int64{int64(x), int64(y)}
	i := 0
	m.In = func() int64 { v := in[i]; i++; return v }
	pulled := false
	m.Out = func(v int64) { pulled = v == 1 }
	m.Run()
	return pulled
}

func part1(program []int64) int {
	pulled := 0
	for x := 0; x < 50; x++ {
		for y := 0; y < 50; y++ {
			if isPulled(program, x, y) {
				pulled++
			}
		}
	}
	return pulled
}

// part2 finds the top-left corner of the smallest 100x100 square fully
// inside the beam and returns X*10000+Y.
//
// The beam is contiguous in every row, and both its left and right
// edges only move right as y grows. So the left edge of each row can
// be found by advancing a pointer that never moves backwards, making
// the whole search O(width of beam) probes instead of O(n^2).
func part2(program []int64) int {
	const size = 100
	const maxRows = 100000

	leftEdge := 0 // left edge of the current row (never decreases)
	for row := 0; row < maxRows; row++ {
		start := leftEdge
		for leftEdge <= 2*row+10 && !isPulled(program, leftEdge, row) {
			leftEdge++
		}
		if leftEdge > 2*row+10 {
			// Empty row: leave the pointer where it was so a later,
			// non-empty row is not skipped.
			leftEdge = start
			continue
		}

		// Treat this row as the bottom row of the square. The left edge
		// is the smallest valid x, so checking the top-right corner is
		// enough to prove all four corners are pulled.
		top := row - (size - 1)
		if top >= 0 && isPulled(program, leftEdge+size-1, top) {
			return leftEdge*10000 + top
		}
	}

	log.Fatal("part2: no square found")
	return 0
}

func main() {
	program, err := intcode.LoadProgram(puzzle.Path("day19.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	fmt.Println("Part1:", part1(program))

	start := time.Now()
	fmt.Println("Part2:", part2(program))
	fmt.Println("Part2 took:", time.Since(start))
}
