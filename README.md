# advent-of-code-2019

Go solutions for [Advent of Code 2019](https://adventofcode.com/2019).

## Layout

- `day01/` … `day25/` — one `package main` per day (`main.go` in each).
- `inputs/` — puzzle inputs (`dayXX.txt`; git-ignored, kept local).
- `internal/puzzle/` — shared helper resolving input files relative to
  the repo root, so days run from any working directory.
- `internal/intcode/` — the one shared Intcode VM (all opcodes, parameter
  modes, relative base, unbounded memory; I/O via callbacks, plus
  step-wise `Feed`/`RunUntilOutput`). Used by days 2, 5, 9, 11, 13, 15,
  17, 19, 21. Days 7, 23, 25 still carry their own VM variants
  (concurrent amplifiers, packet network, interactive adventure).
- `go.mod` / `go.sum` — single Go module (`aoc/2019`).
  Only external dependency is `gonum` (used by day 25).

## Prerequisites

- Go 1.21+

## Run

From the repo root:

```sh
# one day (works from any directory)
go run ./day01

# all days
for d in day*/; do echo "== $d =="; go run ./"$d"; done

# checks
gofmt -l .
go vet ./...
go build ./...
go test ./...
```

Inputs must exist as `inputs/dayXX.txt` (e.g. `inputs/day01.txt`).
Each program locates its input via `internal/puzzle`, which anchors to
the repo root — no need to `cd` into the day's folder.

## Notes

- Line endings are normalized to LF (`.gitattributes`).
- No `go.work` — this is a single-module repo.
