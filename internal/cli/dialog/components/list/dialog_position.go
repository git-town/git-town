package list

// provides the position of the given needle in the given list of entries
func DialogPosition[S comparable](entries Entries[S], needle S) int {
	for e, entry := range entries {
		if entry.Data == needle {
			return e
		}
	}
	return -1
}
