package main
 
import (
        "bufio"
        "fmt"
        "log"
        "os"
        "strconv"
)
 
func main() {
        file, err := os.Open("../data/day01.txt")

        if err != nil {
                log.Fatalf("failed opening file: %s", err)
        }

        scanner := bufio.NewScanner(file)
        scanner.Split(bufio.ScanLines)

        part1 := 0
        part2 := 0

        for scanner.Scan() {
                i, err := strconv.Atoi(scanner.Text())
                if err != nil {
                        // handle error
                        fmt.Println(err)
                        os.Exit(2)
                }

                moduleFuel := int(float64(i/3) - 2)
                part1 += moduleFuel
                part2 += moduleFuel

                for  {
                        moduleFuel = int(float64(moduleFuel/3) - 2)

                        if moduleFuel <= 0 {
                                break
                        }

                        part2 += moduleFuel
                }
        }

        file.Close()

        fmt.Println("Part1:", part1)
        fmt.Println("Part2:", part2)
}