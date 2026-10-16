package puzzle

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRootIsRepoRoot(t *testing.T) {
	root := Root()
	if !filepath.IsAbs(root) {
		t.Fatalf("Root() = %q, want absolute path", root)
	}
	if filepath.Base(root) != "advent-of-code-2019" {
		t.Fatalf("Root() = %q, want repo root", root)
	}
}

func TestPathPointsIntoInputs(t *testing.T) {
	p := Path("day01.txt")
	if filepath.Base(filepath.Dir(p)) != "inputs" {
		t.Fatalf("Path() = %q, want file inside inputs dir", p)
	}
	if !strings.HasSuffix(p, "day01.txt") {
		t.Fatalf("Path() = %q, want day01.txt suffix", p)
	}
}
