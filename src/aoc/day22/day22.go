package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

type operation int

const (
	dealIntoStack     operation = 0
	cut               operation = 1
	dealWithIncrement operation = 2
)

type instruction struct {
	op     operation
	number int
}

func applyDealIntoStack(cards []int64) []int64 {
	for i, j := 0, len(cards)-1; i < j; i, j = i+1, j-1 {
		cards[i], cards[j] = cards[j], cards[i]
	}

	return cards
}

func applyCut(cards []int64, cutNumber int) []int64 {
	var result, tail, head = []int64{}, []int64{}, []int64{}

	if cutNumber > 0 {
		tail = cards[0:cutNumber]
		head = cards[cutNumber:]

	} else {
		head = cards[len(cards)+cutNumber:]
		tail = cards[0 : len(cards)+cutNumber]
	}

	result = append(result, head...)
	result = append(result, tail...)

	return result
}

func applyDealWithIncrement(cards []int64, dealWithAmount int) []int64 {
	var result = make([]int64, len(cards))
	result[0] = cards[0]

	var shuffleIndex = 0
	for i := 1; i < len(cards); i++ {
		shuffleIndex += dealWithAmount
		shuffleIndex = shuffleIndex % len(cards)

		result[shuffleIndex] = cards[i]
	}

	return result
}

func applyShuffling(cards []int64, instructions []instruction) []int64 {
	for _, i := range instructions {
		switch i.op {
		case dealIntoStack:
			cards = applyDealIntoStack(cards)
		case cut:
			cards = applyCut(cards, i.number)
		case dealWithIncrement:
			cards = applyDealWithIncrement(cards, i.number)
		}
	}

	return cards
}

func modInv(n *big.Int, cards *big.Int) *big.Int {
	return big.NewInt(0).Exp(n, big.NewInt(cards.Int64()-2), cards)
}

func applyReverseInstructionShuffling(instructions []instruction, cards *big.Int) (*big.Int, *big.Int) {
	incrementMul := big.NewInt(1)
	offsetDiff := big.NewInt(0)

	for _, i := range instructions {
		switch i.op {
		case dealIntoStack:
			incrementMul.Mul(incrementMul, big.NewInt(-1))
			incrementMul.Mod(incrementMul, cards)
			offsetDiff.Add(offsetDiff, incrementMul)
			offsetDiff.Mod(offsetDiff, cards)
		case cut:
			num := big.NewInt(int64(i.number))
			offsetDiff.Add(offsetDiff, num.Mul(num, incrementMul))
			offsetDiff.Mod(offsetDiff, cards)
		case dealWithIncrement:
			num := big.NewInt(int64(i.number))
			incrementMul.Mul(incrementMul, modInv(num, cards))
			incrementMul.Mod(incrementMul, cards)
		}
	}

	return incrementMul, offsetDiff
}

func getSequence(iterations *big.Int, cards *big.Int, incrementMul *big.Int, offsetDiff *big.Int) (*big.Int, *big.Int) {
	increment := big.NewInt(0)
	increment.Exp(incrementMul, iterations, cards)

	offset := big.NewInt(1)
	one := big.NewInt(1)
	inverse := modInv(one.Sub(one, incrementMul).Mod(one, cards), cards)

	offset.Mul(offset, offsetDiff.Mul(offsetDiff, big.NewInt(1-increment.Int64())).Mul(offsetDiff, inverse))
	offset.Mod(offset, cards)

	return increment, offset
}

func buildInstructions() []instruction {
	file, err := os.Open("../../data/day22.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var instructions = []instruction{}

	for scanner.Scan() {

		line := scanner.Text()
		var i = instruction{}

		if strings.Contains(line, "deal into") {
			i = instruction{dealIntoStack, -1}

		} else if strings.Contains(line, "cut") {
			cutAmount, _ := strconv.Atoi(strings.Split(line, " ")[1])
			i = instruction{cut, cutAmount}
		} else {
			dealWithAmount, _ := strconv.Atoi(strings.Split(line, " ")[3])
			i = instruction{dealWithIncrement, dealWithAmount}
		}

		instructions = append(instructions, i)
	}

	file.Close()
	return instructions
}

func findPositionOfCard2019() int {
	instructions := buildInstructions()

	var cards = []int64{}
	for x := int64(0); x < 10007; x++ {
		cards = append(cards, x)
	}

	cards = applyShuffling(cards, instructions)
	position := 0

	for i, x := range cards {
		if x == 2019 {
			position = i
			break
		}
	}

	return position
}

// https://www.reddit.com/r/adventofcode/comments/ee0rqi/comment/fbnkaju/
// https://github.com/mcpower/adventofcode/blob/501b66084b0060e0375fc3d78460fb549bc7dfab/2019/22/a-improved.py
func findNumberOnCardAtPosition2020() *big.Int {
	instructions := buildInstructions()
	cards := big.NewInt(int64(119315717514047))
	iterations := big.NewInt(int64(101741582076661))

	incrementMul, offsetDiff := applyReverseInstructionShuffling(instructions, cards)
	increment, offset := getSequence(iterations, cards, incrementMul, offsetDiff)

	increment.Mul(increment, big.NewInt(2020))
	offset.Add(offset, increment)
	offset.Mod(offset, cards)
	return offset
}

func main() {
	fmt.Println("Part1:", findPositionOfCard2019())
	fmt.Println("Part2:", findNumberOnCardAtPosition2020())
}
