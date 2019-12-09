package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func decodePixel(pixel string) string {
	decoded := ""
	for i := 0; i < len(pixel); i++ {
		currPixel := pixel[i]

		if (currPixel == '0' || currPixel == '1') && decoded == "" {
			decoded = string(currPixel)
		}
	}

	return decoded
}

func main() {
	file, err := os.Open("../../data/day08.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	imageData := ""

	layerWidth := 25
	layerHeight := 6

	layers := []string{}

	for scanner.Scan() {
		imageData = scanner.Text()
	}

	layerArea := layerWidth * layerHeight

	currentLayer := ""

	for i := 0; i <= len(imageData); i++ {
		if i > 0 && i%layerArea == 0 {
			layers = append(layers, currentLayer)
			currentLayer = ""
		}

		if i == len(imageData) {
			break
		}

		currentLayer += string(imageData[i])
	}

	zeroCounter := 0
	lowestZeroes := math.MaxInt32
	lowestZeroesIndex := 0

	for i, l := range layers {
		for x := 0; x < len(l); x++ {
			if l[x] == '0' {
				zeroCounter++
			}
		}

		if zeroCounter < lowestZeroes {
			lowestZeroes = zeroCounter
			lowestZeroesIndex = i
		}

		zeroCounter = 0
	}

	layerWithLowestZeroes := layers[lowestZeroesIndex]
	countOnes := 0
	countTwos := 0

	for j := 0; j < len(layerWithLowestZeroes); j++ {
		if layerWithLowestZeroes[j] == '1' {
			countOnes++
		} else if layerWithLowestZeroes[j] == '2' {
			countTwos++
		}
	}

	fmt.Println("Part1:", countOnes*countTwos)

	// 0 - black
	// 1 - white
	// 2 - transparent

	decodedImage := ""

	for k := 0; k < layerArea; k++ {
		pixel := ""
		for _, layer := range layers {
			pixel += string(layer[k])
		}
		decodedImage += decodePixel(pixel)
	}

	normalisedImage := ""

	for x := 0; x <= len(decodedImage); x++ {
		if x > 0 && x%layerWidth == 0 {
			normalisedImage += "\n"
		}

		if x == len(decodedImage) {
			break
		}

		normalisedImage += string(decodedImage[x])
	}

	normalisedImage = strings.Replace(normalisedImage, "1", "#", -1)
	normalisedImage = strings.Replace(normalisedImage, "0", ".", -1)

	fmt.Println("Part2:")
	fmt.Println(normalisedImage)

	file.Close()
}
