package format

import (
	"strings"

	"github.com/git-town/git-town/v10/src/config"
)

// BranchLineage provides printable formatting of the given branch lineage.
func BranchLineage(lineage config.Lineage) string {
	roots := lineage.Roots()
	trees := make([]string, len(roots))
	for r, root := range roots {
		trees[r] = BranchTree(root, lineage)
	}
	return strings.Join(trees, "\n\n")
}
