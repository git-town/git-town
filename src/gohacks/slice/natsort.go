package slice

import (
	"fmt"
	"sort"
	"strconv"
	"unicode"
)

func NatSort[T fmt.Stringer](elements []T) []T {
	stringerSlice := make(Stringers, len(elements))
	for e, element := range elements {
		stringerSlice[e] = element
	}
	sort.Sort(stringerSlice)
	result := make([]T, len(stringerSlice))
	for s, stringer := range stringerSlice {
		result[s] = stringer.(T)
	}
	return result
}

type Stringers []fmt.Stringer

func (self Stringers) Len() int {
	return len(self)
}

func (self Stringers) Less(a, b int) bool {
	return naturalLess(self[a].String(), self[b].String())
}

func (self Stringers) Swap(a, b int) {
	self[a], self[b] = self[b], self[a]
}

func naturalLess(text1, text2 string) bool {
	index1, index2 := 0, 0
	for index1 < len(text1) && index2 < len(text2) {
		var part1, part2 string
		if unicode.IsDigit(rune(text1[index1])) {
			part1, index1 = extractNumber(text1, index1)
		} else {
			part1, index1 = extractNonNumber(text1, index1)
		}
		if unicode.IsDigit(rune(text2[index2])) {
			part2, index2 = extractNumber(text2, index2)
		} else {
			part2, index2 = extractNonNumber(text2, index2)
		}
		if part1 != part2 {
			if unicode.IsDigit(rune(part1[0])) && unicode.IsDigit(rune(part2[0])) {
				// compare numbers by their numeric value
				int1, _ := strconv.Atoi(part1)
				int2, _ := strconv.Atoi(part2)
				return int1 < int2
			}
			// compare non-numbers lexicographically
			return part1 < part2
		}
	}
	// the strings are equal up to the end of one of them
	return len(text1) < len(text2)
}

func extractNumber(text string, index int) (number string, nextIndex int) {
	for nextIndex = index; nextIndex < len(text) && unicode.IsDigit(rune(text[nextIndex])); nextIndex++ {
	}
	return text[index:nextIndex], nextIndex
}

func extractNonNumber(text string, index int) (nonNumber string, nextIndex int) {
	for nextIndex = index; nextIndex < len(text) && !unicode.IsDigit(rune(text[nextIndex])); nextIndex++ {
	}
	return text[index:nextIndex], nextIndex
}
