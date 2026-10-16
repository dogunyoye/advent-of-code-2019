package intcode

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func runCollect(t *testing.T, program []int64, inputs ...int64) []int64 {
	t.Helper()
	m := New(program)
	for _, in := range inputs {
		m.Feed(in)
	}
	var outputs []int64
	m.Out = func(v int64) { outputs = append(outputs, v) }
	m.Run()
	if !m.Halted {
		t.Fatal("machine did not halt")
	}
	return outputs
}

func TestQuine(t *testing.T) {
	program := []int64{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99}
	if got := runCollect(t, program); !reflect.DeepEqual(got, program) {
		t.Fatalf("quine mismatch:\n got %v\nwant %v", got, program)
	}
}

func TestLargeNumber(t *testing.T) {
	got := runCollect(t, []int64{1102, 34915192, 34915192, 7, 4, 7, 99, 0})
	if len(got) != 1 || len(toString(got[0])) != 16 {
		t.Fatalf("want 16-digit number, got %v", got)
	}
}

func TestImmediateOutput(t *testing.T) {
	got := runCollect(t, []int64{104, 1125899906842624, 99})
	if len(got) != 1 || got[0] != 1125899906842624 {
		t.Fatalf("got %v", got)
	}
}

func TestRelativeIOEcho(t *testing.T) {
	// NOTE: the data cell (rel+0 = 10) must not overlap the code:
	// with rel=5 the input write would clobber the 204 instruction's
	// own parameter cell.
	got := runCollect(t, []int64{109, 10, 203, 0, 204, 0, 99}, 42)
	if len(got) != 1 || got[0] != 42 {
		t.Fatalf("got %v", got)
	}
}

func TestParameterModes(t *testing.T) {
	m := New([]int64{1002, 4, 3, 4, 33})
	m.Run()
	if m.Mem(4) != 99 {
		t.Fatalf("Mem(4) = %d, want 99", m.Mem(4))
	}

	m = New([]int64{1101, 100, -1, 4, 0})
	m.Run()
	if m.Mem(4) != 99 {
		t.Fatalf("Mem(4) = %d, want 99", m.Mem(4))
	}
}

func TestJumpAndCompare(t *testing.T) {
	program := []int64{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}
	if got := runCollect(t, program, 8); len(got) != 1 || got[0] != 1 {
		t.Fatalf("input 8: got %v, want [1]", got)
	}
	if got := runCollect(t, program, 7); len(got) != 1 || got[0] != 0 {
		t.Fatalf("input 7: got %v, want [0]", got)
	}
}

func TestRunUntilOutput(t *testing.T) {
	m := New([]int64{3, 0, 4, 0, 3, 0, 4, 0, 99})
	m.Feed(7)
	v, halted := m.RunUntilOutput()
	if halted || v != 7 {
		t.Fatalf("first = (%d, %v), want (7, false)", v, halted)
	}
	m.Feed(8)
	v, halted = m.RunUntilOutput()
	if halted || v != 8 {
		t.Fatalf("second = (%d, %v), want (8, false)", v, halted)
	}
	if _, halted = m.RunUntilOutput(); !halted {
		t.Fatal("want halted after outputs consumed")
	}
}

func TestStepSequence(t *testing.T) {
	m := New([]int64{3, 0, 4, 0, 99})

	if ev := m.Step(); ev.Kind != NeedInput {
		t.Fatalf("step 1 = %+v, want NeedInput", ev)
	}
	// Nothing executed: still parked on the input instruction.
	if m.IP() != 0 {
		t.Fatalf("IP = %d, want 0", m.IP())
	}

	m.Feed(7)
	if ev := m.Step(); ev.Kind != Continue {
		t.Fatalf("step 2 = %+v, want Continue", ev)
	}
	if ev := m.Step(); ev.Kind != Output || ev.Value != 7 {
		t.Fatalf("step 3 = %+v, want Output 7", ev)
	}
	if ev := m.Step(); ev.Kind != Halted {
		t.Fatalf("step 4 = %+v, want Halted", ev)
	}
	if !m.Halted {
		t.Fatal("want Halted flag set")
	}
	// Stepping a halted machine stays halted.
	if ev := m.Step(); ev.Kind != Halted {
		t.Fatalf("step 5 = %+v, want Halted", ev)
	}
}

func TestStepUsesInCallback(t *testing.T) {
	m := New([]int64{3, 0, 4, 0, 99})
	m.In = func() int64 { return 9 }
	if ev := m.Step(); ev.Kind != Continue {
		t.Fatalf("step 1 = %+v, want Continue", ev)
	}
	if ev := m.Step(); ev.Kind != Output || ev.Value != 9 {
		t.Fatalf("step 2 = %+v, want Output 9", ev)
	}
}

func TestIPPointsAtNextInstruction(t *testing.T) {
	m := New([]int64{1101, 100, -1, 4, 0})
	if m.IP() != 0 {
		t.Fatalf("IP = %d, want 0", m.IP())
	}
	m.Run()
	if m.IP() != 4 {
		t.Fatalf("IP = %d, want 4 (halt address)", m.IP())
	}
}

func TestNewCopiesProgram(t *testing.T) {
	template := []int64{1, 0, 0, 0, 99}
	m := New(template)
	m.Run()
	if m.Mem(0) != 2 {
		t.Fatalf("Mem(0) = %d, want 2", m.Mem(0))
	}
	if template[0] != 1 {
		t.Fatalf("template mutated: %v", template)
	}
}

func TestMemoryGrowsWithZeros(t *testing.T) {
	m := New([]int64{99})
	if m.Mem(1000) != 0 {
		t.Fatalf("Mem(1000) = %d, want 0", m.Mem(1000))
	}
	m.SetMem(1000, 5)
	if m.Mem(1000) != 5 {
		t.Fatalf("Mem(1000) = %d, want 5", m.Mem(1000))
	}
}

func TestLoadProgram(t *testing.T) {
	path := filepath.Join(t.TempDir(), "prog.txt")
	if err := os.WriteFile(path, []byte("1,0,0,0,99\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	program, err := LoadProgram(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(program, []int64{1, 0, 0, 0, 99}) {
		t.Fatalf("got %v", program)
	}
}

func toString(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		return "-" + string(digits)
	}
	return string(digits)
}
