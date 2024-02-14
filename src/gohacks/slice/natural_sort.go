package slice

import (
	"fmt"
	"sort"
	"strconv"
	"unicode"
)

// sorts the given elements in natural sort order (https://en.wikipedia.org/wiki/Natural_sort_order)
func NaturalSort[T fmt.Stringer](elements []T) []T {
	stringers := newStringers(elements)
	sort.Sort(stringers)
	return stringers.unwrap()
}

func extractNonNumber(text string, index int) (nonNumber string, nextIndex int) {
	for nextIndex = index; nextIndex < len(text) && !unicode.IsDigit(rune(text[nextIndex])); nextIndex++ { //revive:disable-line:empty-block
	}
	return text[index:nextIndex], nextIndex
}

func extractNumber(text string, index int) (number string, nextIndex int) {
	for nextIndex = index; nextIndex < len(text) && unicode.IsDigit(rune(text[nextIndex])); nextIndex++ { //revive:disable-line:empty-block
	}
	return text[index:nextIndex], nextIndex
}

// indicates whether text1 < text2 according to natural sort order
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

type stringers[T fmt.Stringer] []fmt.Stringer

func newStringers[T fmt.Stringer](elements []T) stringers[T] {
	stringers := make(stringers[T], len(elements))
	for e, element := range elements {
		stringers[e] = element
	}
	return stringers
}

func (self stringers[T]) Len() int {
	return len(self)
}

func (self stringers[T]) Less(a, b int) bool {
	return naturalLess(self[a].String(), self[b].String())
}

func (self stringers[T]) Swap(a, b int) {
	self[a], self[b] = self[b], self[a]
}

func (self stringers[T]) unwrap() []T {
	result := make([]T, len(self))
	for s, stringer := range self {
		result[s] = stringer.(T) //nolint:forcetypeassert  // we are sure this is T here because stringers consists of T
	}
	return result
}
