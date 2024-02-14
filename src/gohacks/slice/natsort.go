package slice

import (
	"fmt"
	"sort"
	"strconv"
	"unicode"
)

type StringerSlice []fmt.Stringer

func (s StringerSlice) Len() int {
	return len(s)
}

func (s StringerSlice) Less(a, b int) bool {
	return naturalLess(s[a].String(), s[b].String())
}

func (s StringerSlice) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

func NatSort[T fmt.Stringer](elements []T) {
	s := StringerSlice{}
}

func SortStringers(s StringerSlice) {
	sort.Sort(s)
}

func naturalLess(a, b string) bool {
	ai, bi := 0, 0
	for ai < len(a) && bi < len(b) {
		var an, bn string
		if unicode.IsDigit(rune(a[ai])) {
			an, ai = extractNumber(a, ai)
		} else {
			an, ai = extractNonNumber(a, ai)
		}
		if unicode.IsDigit(rune(b[bi])) {
			bn, bi = extractNumber(b, bi)
		} else {
			bn, bi = extractNonNumber(b, bi)
		}
		if an != bn {
			if unicode.IsDigit(rune(an[0])) && unicode.IsDigit(rune(bn[0])) {
				// Compare numbers by their numeric value.
				anInt, _ := strconv.Atoi(an)
				bnInt, _ := strconv.Atoi(bn)
				return anInt < bnInt
			}
			// Compare non-numbers lexicographically.
			return an < bn
		}
	}
	// The strings are equal up to the end of one of them.
	return len(a) < len(b)
}

func extractNumber(s string, i int) (number string, next int) {
	for next = i; next < len(s) && unicode.IsDigit(rune(s[next])); next++ {
	}
	return s[i:next], next
}

func extractNonNumber(s string, i int) (nonNumber string, next int) {
	for next = i; next < len(s) && !unicode.IsDigit(rune(s[next])); next++ {
	}
	return s[i:next], next
}
