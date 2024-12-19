// Package dialog provides high-level screens through which the user can enter data into Git Town.
package dialog

import "github.com/git-town/git-town/v17/internal/cli/dialog/components/list"

func DialogPosition[S comparable](entries list.Entries[S], needle S) int {
	for e, entry := range entries {
		if entry.Data == needle {
			return e
		}
	}
	return -1
}
