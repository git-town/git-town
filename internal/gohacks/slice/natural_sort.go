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
	cutter1 := newCutter(text1)
	cutter2 := newCutter(text2)
	for cutter1.hasMoreParts() && cutter2.hasMoreParts() {
		part1 := cutter1.nextPart()
		part2 := cutter2.nextPart()
		if part1 != part2 {
			if part1.isNumber() && part2.isNumber() {
				// compare numbers by their numeric value
				return part1.toNumber() < part2.toNumber()
			}
			// compare non-numbers lexicographically
			return part1 < part2
		}
	}
	// the strings are equal up to the end of one of them
	return len(text1) < len(text2)
}

// cuts given text into parts of consecutive numbers and non-numbers
type cutter struct {
	content string // the content to cut into parts
	index   int    // where we are in the content right now
}

func newCutter(content string) cutter {
	return cutter{
		content: content,
		index:   0,
	}
}

// indicates whether the given index is inside the content this cutter disects
func (c cutter) hasContentAt(index int) bool {
	return index < len(c.content)
}

// indicates whether this cutter can yield more parts
func (c cutter) hasMoreParts() bool {
	return c.index < len(c.content)
}

// indicates whether the rune at the given index is a number
func (c cutter) isDigitAt(index int) bool {
	return unicode.IsDigit(rune(c.content[index]))
}

// provides the next part of the content that this cutter disects
func (c *cutter) nextPart() part {
	index := c.index
	lookingForDigits := c.isDigitAt(index)
	for c.hasContentAt(index) && c.isDigitAt(index) == lookingForDigits {
		index++
	}
	result := part(c.content[c.index:index])
	c.index = index
	return result
}

// a multi-rune part cut from text, either all numbers or all non-numbers
type part string

// indicates whether this part contains a number
func (part part) isNumber() bool {
	return unicode.IsDigit(rune(part[0]))
}

// assuming this block contains numbers, provides the resulting numeric value
func (part part) toNumber() int {
	result, _ := strconv.Atoi(string(part))
	return result
}

// wraps the given []fmt.Stringer with the sort.Interface methods so that we can sort it using the stdlib
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
