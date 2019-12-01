package main
 
import (
        "bufio"
        "fmt"
        "log"
        "os"
        "strconv"
        "math"
)
 
func main() {
        file, err := os.Open("../data/day01.txt")

        if err != nil {
                log.Fatalf("failed opening file: %s", err)
        }

        scanner := bufio.NewScanner(file)
        scanner.Split(bufio.ScanLines)
        var txtlines []string

        for scanner.Scan() {
                txtlines = append(txtlines, scanner.Text())
        }

        file.Close()

        part1 := 0
        part2 := 0

        for _, eachline := range txtlines {
                i, err := strconv.Atoi(eachline)
                if err != nil {
                        // handle error
                        fmt.Println(err)
                        os.Exit(2)
                }

                moduleFuel := int(math.Floor(float64(i/3)) - 2)
                part1 += moduleFuel
                part2 += moduleFuel

                for  {
                        moduleFuel = int(math.Floor(float64(moduleFuel/3)) - 2)

                        if moduleFuel <= 0 {
                                break
                        }

                        part2 += moduleFuel
                }
        }

        fmt.Println("Part1:", part1)
        fmt.Println("Part2:", part2)
}