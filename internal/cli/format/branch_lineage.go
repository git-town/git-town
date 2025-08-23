package format

import (
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const indent = "  "

// BranchLineage provides printable formatting of the given branch lineage.
func BranchLineage(lineage configdomain.Lineage) string {
	result := strings.Builder{}
	for _, root := range lineage.Roots() {
		branchTree(branchTreeArgs{
			branch:      root,
			builder:     NewMutable(&result),
			indentLevel: 0,
			lineage:     lineage,
		})
	}
	return result.String()
}

type branchTreeArgs struct {
	branch      gitdomain.LocalBranchName
	builder     Mutable[strings.Builder]
	indentLevel int
	lineage     configdomain.Lineage
}

// branchTree provids a printable version of the given branch tree.
func branchTree(args branchTreeArgs) {
	for range args.indentLevel {
		args.builder.Value.WriteString(indent)
	}
	args.builder.Value.WriteString(args.branch.String())
	for _, child := range args.lineage.Children(args.branch) {
		args.builder.Value.WriteString("\n")
		branchTree(branchTreeArgs{
			branch:      child,
			builder:     args.builder,
			indentLevel: args.indentLevel + 1,
			lineage:     args.lineage,
		})
	}
}
