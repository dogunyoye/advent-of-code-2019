package main

import (
	"bufio"
	"fmt"
	"log"
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

func main() {
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

	//fmt.Println(instructions)

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

	fmt.Println("Part1:", position)
}
