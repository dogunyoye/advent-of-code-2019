package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"aoc/2019/internal/intcode"
	"aoc/2019/internal/puzzle"
)

type amplifier struct {
	code   string
	input  chan int64
	output chan int64
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

func runAmplifier(program []int64, phase, signal int64) int64 {
	var last int64
	var numOfInputs = 0

	m := intcode.New(program)
	m.In = func() int64 {
		numOfInputs++

		// on every second input request, we take the previous amplifiers result
		if numOfInputs%2 == 0 {
			return signal
		}
		return phase
	}
	m.Out = func(v int64) { last = v }
	m.Run()
	return last
}

func runAmplifierFeedback(phase int64, amp amplifier, wg *sync.WaitGroup, completed *int32, result *int64, program []int64) {
	var numOfInputs = 0

	var ampAInitialised = false

	m := intcode.New(program)
	m.In = func() int64 {
		numOfInputs++

		if numOfInputs == 1 {
			return phase
		}
		if amp.code == "A" && !ampAInitialised {
			ampAInitialised = true
			return 0
		}
		return <-amp.input
	}
	m.Out = func(v int64) {
		if amp.code == "E" {
			*result = v
			// super hacky sleep to ensure the last amplifier doesn't
			// race ahead of the previous amplifier (D) to ensure a
			// correct value for completed
			time.Sleep(1 * time.Microsecond)
		}

		if atomic.LoadInt32(completed) > 1 {
			amp.output <- v
		}
	}
	m.Run()

	atomic.AddInt32(completed, -1)
	wg.Done()
}

func main() {

	program, err := intcode.LoadProgram(puzzle.Path("day07.txt"))
	if err != nil {
		log.Fatalf("failed loading program: %s", err)
	}

	combinations := generateCombinations([]string{"0", "1", "2", "3", "4"})

	var max int64
	maxCombo := ""

	for _, combo := range combinations {
		var signal int64
		for i := 0; i < len(combo); i++ {
			phase, _ := strconv.Atoi(string(combo[i]))
			signal = runAmplifier(program, int64(phase), signal)
		}

		if signal > max {
			max = signal
			maxCombo = combo
		}
	}

	fmt.Println("Part1:", max, maxCombo)
	max = 0
	maxCombo = ""

	combinationsPart2 := generateCombinations([]string{"5", "6", "7", "8", "9"})

	var aChannel = make(chan int64)
	var bChannel = make(chan int64)
	var cChannel = make(chan int64)
	var dChannel = make(chan int64)
	var eChannel = make(chan int64)

	var ampA = amplifier{"A", aChannel, bChannel}
	var ampB = amplifier{"B", bChannel, cChannel}
	var ampC = amplifier{"C", cChannel, dChannel}
	var ampD = amplifier{"D", dChannel, eChannel}
	var ampE = amplifier{"E", eChannel, aChannel}

	var amplifiers = [5]amplifier{ampA, ampB, ampC, ampD, ampE}

	for _, comboPart2 := range combinationsPart2 {
		var completed = int32(5)
		var wg sync.WaitGroup
		var result int64

		for i := 0; i < len(comboPart2); i++ {
			phase, _ := strconv.Atoi(string(comboPart2[i]))
			wg.Add(1)
			go runAmplifierFeedback(int64(phase), amplifiers[i], &wg, &completed, &result, program)
		}

		// block here until last amplifer halts
		wg.Wait()

		if result > max {
			max = result
			maxCombo = comboPart2
		}
	}

	fmt.Println("Part2:", max, maxCombo)
}
