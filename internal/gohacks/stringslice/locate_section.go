package stringslice

// LocateSection locates the line number where the given lines start in the given lines.
func LocateSection(lines, section []string) (int, bool) {
	for i := 0; i <= len(lines)-len(section); i++ {
		if Matches(lines[i:i+len(section)], section) {
			return i, true
		}
	}
	return -1, false
}
