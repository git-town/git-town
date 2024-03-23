package format

import (
	"strings"

	"github.com/git-town/git-town/v13/src/config/configdomain"
)

// BranchLineage provides printable formatting of the given branch lineage.
func BranchLineage(lineage configdomain.Lineage) string {
	roots := lineage.Roots()
	trees := make([]string, len(roots))
	for r, root := range roots {
		trees[r] = BranchTree(root, lineage)
	}
	return strings.Join(trees, "\n\n")
}
