package slice

import (
	"fmt"
	"sort"
	"strconv"
	"unicode"
)

// sorts the given elements in natural sort order (https://en.wikipedia.org/wiki/Natural_sort_order)
func NaturalSort[T fmt.Stringer](list []T) []T {
	sortableList := newSortable(list)
	sort.Sort(sortableList)
	return sortableList
}

// indicates whether text1 < text2 according to natural sort order
func naturalLess(text1, text2 string) bool {
	cursor1 := newCursor(text1)
	cursor2 := newCursor(text2)
	for cursor1.hasMore() && cursor2.hasMore() {
		part1 := cursor1.nextPart()
		part2 := cursor2.nextPart()
		if part1 != part2 {
			if part1.isNumber() && part2.isNumber() {
				// compare numbers by their numeric value
				return part1.asNumber() < part2.asNumber()
			}
			// compare non-numbers lexicographically
			return part1 < part2
		}
	}
	// the strings are equal up to the end of one of them
	return len(text1) < len(text2)
}

type cursor struct {
	index int
	text  string
}

func (cursor cursor) hasMore() bool {
	return cursor.index < len(cursor.text)
}

func (cursor *cursor) nextPart() part {
	var endIndex int
	if unicode.IsDigit(rune(cursor.text[cursor.index])) {
		for endIndex = cursor.index; endIndex < len(cursor.text) && unicode.IsDigit(rune(cursor.text[endIndex])); endIndex++ { //revive:disable-line:empty-block
		}
	} else {
		for endIndex = cursor.index; endIndex < len(cursor.text) && !unicode.IsDigit(rune(cursor.text[endIndex])); endIndex++ { //revive:disable-line:empty-block
		}
	}
	result := part(cursor.text[cursor.index:endIndex])
	cursor.index = endIndex
	return result
}

// a part of a search string, represents either a multi-digit number or a text block
type part string

func (part part) isNumber() bool {
	return unicode.IsDigit(rune(part[0]))
}

func (part part) asNumber() int {
	result, _ := strconv.Atoi(string(part))
	return result
}

func newCursor(text string) cursor {
	return cursor{
		index: 0,
		text:  text,
	}
}

// wraps the given []fmt.Stringer with methods that allow sorting it using the stdlib
type sortable[T fmt.Stringer] []T

func newSortable[T fmt.Stringer](elements []T) sortable[T] {
	sortable := make(sortable[T], len(elements))
	copy(sortable, elements)
	return sortable
}

func (self sortable[T]) Len() int {
	return len(self)
}

func (self sortable[T]) Less(a, b int) bool {
	return naturalLess(self[a].String(), self[b].String())
}

func (self sortable[T]) Swap(a, b int) {
	self[a], self[b] = self[b], self[a]
}
