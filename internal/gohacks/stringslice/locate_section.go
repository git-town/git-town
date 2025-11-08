package stringslice

// LocateSection provides the line number where the given section starts within the given lines.
func LocateSection(lines, section []string) (int, bool) {
	sectionLength := len(section)
	for i := 0; i <= len(lines)-sectionLength; i++ {
		if EqualIgnoreWhitespace(lines[i:i+sectionLength], section) {
			return i, true
		}
	}
	return -1, false
}
