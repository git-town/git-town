package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/gohacks/slice"
)

type stringer struct {
	s string
}

func (s stringer) String() string {
	return s.s
}

func TestSortStringers(t *testing.T) {
	tests := []struct {
		name     string
		input    slice.StringerSlice
		expected slice.StringerSlice
	}{
		{
			name:     "empty slice",
			input:    slice.StringerSlice{},
			expected: slice.StringerSlice{},
		},
		{
			name:     "single element",
			input:    slice.StringerSlice{stringer{"a"}},
			expected: slice.StringerSlice{stringer{"a"}},
		},
		{
			name:     "multiple elements",
			input:    slice.StringerSlice{stringer{"b10"}, stringer{"b2"}, stringer{"a1"}},
			expected: slice.StringerSlice{stringer{"a1"}, stringer{"b2"}, stringer{"b10"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			slice.SortStringers(test.input)
			for i, v := range test.input {
				if v.String() != test.expected[i].String() {
					t.Errorf("Expected %v, but got %v", test.expected, test.input)
				}
			}
		})
	}
}
