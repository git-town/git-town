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

// indicates whether this cutter can yield more parts
func (c cutter) hasMoreParts() bool {
	return c.index < len(c.content)
}

// provides the next part of the content that this cutter disects
func (c *cutter) nextPart() part {
	endIndex := c.index
	if unicode.IsDigit(rune(c.content[c.index])) {
		for ; endIndex < len(c.content) && unicode.IsDigit(rune(c.content[endIndex])); endIndex++ { //revive:disable-line:empty-block
		}
	} else {
		for ; endIndex < len(c.content) && !unicode.IsDigit(rune(c.content[endIndex])); endIndex++ { //revive:disable-line:empty-block
		}
	}
	result := part(c.content[c.index:endIndex])
	c.index = endIndex
	return result
}

// a part cut from text, either a multi-digit number or a block of non-numbers
type part string

// indicates whether this part is a block of numbers
func (part part) isNumber() bool {
	return unicode.IsDigit(rune(part[0]))
}

// assuming this block contains of numbers, provides their numeric value
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
