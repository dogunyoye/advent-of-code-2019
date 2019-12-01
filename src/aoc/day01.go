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

        requirement := 0
        var moduleFuels []int

        for _, eachline := range txtlines {
                i, err := strconv.Atoi(eachline)
                if err != nil {
                        // handle error
                        fmt.Println(err)
                        os.Exit(2)
                }

                moduleFuel := math.Floor(float64(i/3)) - 2
                requirement += int(moduleFuel)

                moduleFuels = append(moduleFuels, int(moduleFuel))
        }

        fmt.Println("Part1:", requirement)

        total := 0

        for _, fuel := range moduleFuels {
                currentFuel := fuel
                total += currentFuel

                for  {
                        currentFuel = int(math.Floor(float64(currentFuel/3)) - 2)

                        if int(currentFuel) <= 0 {
                                break
                        }

                        total += currentFuel
                }
        }

        fmt.Println("Part2:", total)
}