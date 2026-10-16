// Package puzzle resolves puzzle-input files relative to the repository
// root, so day programs run correctly regardless of the working directory
// they are started from (e.g. `go run ./day01` from the repo root).
package puzzle

import (
	"path/filepath"
	"runtime"
)

// Root returns the absolute path of the repository root.
func Root() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("puzzle: cannot determine source file location")
	}
	// file is <root>/internal/puzzle/puzzle.go.
	return filepath.Dir(filepath.Dir(filepath.Dir(file)))
}

// Path returns the absolute path of the named file under the inputs
// directory (e.g. Path("day01.txt")).
func Path(name string) string {
	return filepath.Join(Root(), "inputs", name)
}
