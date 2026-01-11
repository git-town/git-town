package stringslice

import . "github.com/git-town/git-town/v22/pkg/prelude"

// LocateSection provides the line number where the given section starts within the given lines.
func LocateSection(lines, section []string) Option[int] {
	sectionLength := len(section)
	for i := 0; i <= len(lines)-sectionLength; i++ {
		if EqualIgnoreWhitespace(lines[i:i+sectionLength], section) {
			return Some(i)
		}
	}
	return None[int]()
}
