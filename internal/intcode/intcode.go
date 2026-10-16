// Package intcode implements the Advent of Code 2019 Intcode virtual
// machine: 64-bit words, unbounded zero-initialised memory, opcodes
// 1, 2, 3, 4, 5, 6, 7, 8, 9 and 99, and parameter modes 0 (position),
// 1 (immediate) and 2 (relative).
//
// Input and output are decoupled from the machine via callbacks so every
// day's driver (single values, scripted streams, sensors, networks) can
// plug in without changing the core:
//
//	m := intcode.New(program)
//	m.In = func() int64 { return nextSensorReading() }
//	m.Out = func(v int64) { handleOutput(v) }
//	m.Run()
//
// For step-wise drivers, Feed queues one input value and RunUntilOutput
// runs until the next output (or halt).
package intcode

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Machine is a single Intcode VM instance. Create one with New; the
// zero value is not usable. A Machine is not safe for concurrent use;
// run each computer on its own Machine.
type Machine struct {
	mem []int64
	ip  int64
	rel int64

	// In supplies a value each time opcode 3 executes. When nil,
	// values are taken from the feed queue (see Feed); executing
	// opcode 3 with neither set panics.
	In func() int64
	// Out receives every opcode 4 value. It may be nil, in which case
	// outputs are dropped (RunUntilOutput still reports them).
	Out func(int64)

	queue []int64

	// Halted is true once opcode 99 has executed.
	Halted bool
}

// LoadProgram reads a comma-separated Intcode program from path.
func LoadProgram(path string) ([]int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var text string
	for scanner.Scan() {
		text = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var program []int64
	for _, field := range strings.Split(text, ",") {
		n, err := strconv.ParseInt(strings.TrimSpace(field), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("intcode: bad value %q: %w", field, err)
		}
		program = append(program, n)
	}
	return program, nil
}

// New returns a Machine holding its own copy of program, so the
// caller's slice is never mutated and can seed many Machines.
func New(program []int64) *Machine {
	mem := make([]int64, len(program))
	copy(mem, program)
	return &Machine{mem: mem}
}

// Feed queues an input value consumed by the next opcode 3 when In is
// nil. Useful for step-wise drivers.
func (m *Machine) Feed(v int64) {
	m.queue = append(m.queue, v)
}

// Mem reads memory (out-of-range reads yield 0).
func (m *Machine) Mem(addr int64) int64 {
	return m.get(addr)
}

// SetMem writes memory, growing it as needed.
func (m *Machine) SetMem(addr, value int64) {
	m.set(addr, value)
}

func (m *Machine) get(addr int64) int64 {
	if addr < 0 {
		panic(fmt.Sprintf("intcode: negative address %d", addr))
	}
	if addr >= int64(len(m.mem)) {
		return 0
	}
	return m.mem[addr]
}

func (m *Machine) set(addr, value int64) {
	if addr < 0 {
		panic(fmt.Sprintf("intcode: negative address %d", addr))
	}
	for addr >= int64(len(m.mem)) {
		m.mem = append(m.mem, 0)
	}
	m.mem[addr] = value
}

// EventKind classifies the outcome of a single Step.
type EventKind int

const (
	// Continue means a silent instruction executed (or input consumed).
	Continue EventKind = iota
	// NeedInput means opcode 3 was reached with no input available:
	// In is nil and the feed queue is empty. Nothing executed;
	// Feed a value (or set In) and Step again.
	NeedInput
	// Output means opcode 4 executed; Event.Value holds the value
	// (also delivered to Out when set).
	Output
	// Halted means opcode 99 executed (or the machine already halted).
	Halted
)

// Event is the outcome of a single Step.
type Event struct {
	Kind  EventKind
	Value int64 // valid when Kind == Output
}

// IP returns the address of the next instruction to execute, letting
// schedulers inspect it (e.g. to distinguish a bare opcode 3 word
// from a mode-prefixed one).
func (m *Machine) IP() int64 {
	return m.ip
}

// Step executes a single instruction and reports what happened. It is
// the building block for schedulers that must pause and resume a
// machine (networks, interactive programs).
func (m *Machine) Step() Event {
	if m.Halted {
		return Event{Kind: Halted}
	}
	return m.exec()
}

// Run executes until opcode 99.
func (m *Machine) Run() {
	m.run(false)
}

// RunUntilOutput executes until the next opcode 4 or halt, returning
// the output value (false) or (0, true) on halt. The output is still
// delivered to Out when set.
func (m *Machine) RunUntilOutput() (int64, bool) {
	return m.run(true)
}

func (m *Machine) run(stopOnOutput bool) (int64, bool) {
	for !m.Halted {
		switch ev := m.exec(); ev.Kind {
		case Halted:
			return 0, true
		case Output:
			if stopOnOutput {
				return ev.Value, false
			}
		case NeedInput:
			panic("intcode: opcode 3 with In unset and empty feed queue")
		case Continue:
		}
	}
	return 0, true
}

// exec executes the instruction at the instruction pointer, returning
// its event. A NeedInput event advances nothing.
func (m *Machine) exec() Event {
	op := m.get(m.ip)
	code := op % 100
	var mode [3]int64
	for i, d := 0, op/100; i < 3; i, d = i+1, d/10 {
		mode[i] = d % 10
	}

	switch code {
	case 1:
		m.set(m.writeAddr(3, mode[2]), m.readParam(1, mode[0])+m.readParam(2, mode[1]))
		m.ip += 4
	case 2:
		m.set(m.writeAddr(3, mode[2]), m.readParam(1, mode[0])*m.readParam(2, mode[1]))
		m.ip += 4
	case 3:
		var v int64
		if m.In != nil {
			v = m.In()
		} else if len(m.queue) != 0 {
			v = m.queue[0]
			m.queue = m.queue[1:]
		} else {
			return Event{Kind: NeedInput}
		}
		m.set(m.writeAddr(1, mode[0]), v)
		m.ip += 2
	case 4:
		v := m.readParam(1, mode[0])
		m.ip += 2
		if m.Out != nil {
			m.Out(v)
		}
		return Event{Kind: Output, Value: v}
	case 5:
		if m.readParam(1, mode[0]) != 0 {
			m.ip = m.readParam(2, mode[1])
		} else {
			m.ip += 3
		}
	case 6:
		if m.readParam(1, mode[0]) == 0 {
			m.ip = m.readParam(2, mode[1])
		} else {
			m.ip += 3
		}
	case 7:
		var v int64
		if m.readParam(1, mode[0]) < m.readParam(2, mode[1]) {
			v = 1
		}
		m.set(m.writeAddr(3, mode[2]), v)
		m.ip += 4
	case 8:
		var v int64
		if m.readParam(1, mode[0]) == m.readParam(2, mode[1]) {
			v = 1
		}
		m.set(m.writeAddr(3, mode[2]), v)
		m.ip += 4
	case 9:
		m.rel += m.readParam(1, mode[0])
		m.ip += 2
	case 99:
		m.Halted = true
		return Event{Kind: Halted}
	default:
		panic(fmt.Sprintf("intcode: unknown opcode %d at %d", op, m.ip))
	}
	return Event{Kind: Continue}
}

// readParam returns the value of the n-th parameter (1-based) under mode.
func (m *Machine) readParam(n int, mode int64) int64 {
	raw := m.get(m.ip + int64(n))
	switch mode {
	case 0:
		return m.get(raw)
	case 1:
		return raw
	case 2:
		return m.get(m.rel + raw)
	default:
		panic(fmt.Sprintf("intcode: bad read mode %d at %d", mode, m.ip))
	}
}

// writeAddr returns the target address of the n-th parameter (1-based).
func (m *Machine) writeAddr(n int, mode int64) int64 {
	raw := m.get(m.ip + int64(n))
	switch mode {
	case 0:
		return raw
	case 2:
		return m.rel + raw
	default:
		panic(fmt.Sprintf("intcode: bad write mode %d at %d", mode, m.ip))
	}
}
