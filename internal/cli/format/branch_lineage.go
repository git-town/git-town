package format

import (
	"strings"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const indent = "  "

// BranchLineage provides printable formatting of the given branch lineage.
func BranchLineage(lineage configdomain.Lineage, order configdomain.Order) string {
	roots := lineage.Roots()
	if len(roots) == 0 {
		return ""
	}
	result := strings.Builder{}
	for _, root := range roots {
		result.WriteString("\n\n")
		branchTree(branchTreeArgs{
			branch:      root,
			builder:     NewMutable(&result),
			indentLevel: 0,
			lineage:     lineage,
			order:       order,
		})
	}
	return result.String()[2:]
}

type branchTreeArgs struct {
	branch      gitdomain.LocalBranchName
	builder     Mutable[strings.Builder]
	indentLevel int
	lineage     configdomain.Lineage
	order       configdomain.Order
}

// branchTree provids a printable version of the given branch tree.
func branchTree(args branchTreeArgs) {
	for range args.indentLevel {
		args.builder.Value.WriteString(indent)
	}
	args.builder.Value.WriteString(args.branch.String())
	for _, child := range args.lineage.Children(args.branch, args.order) {
		args.builder.Value.WriteString("\n")
		branchTree(branchTreeArgs{
			branch:      child,
			builder:     args.builder,
			indentLevel: args.indentLevel + 1,
			lineage:     args.lineage,
			order:       args.order,
		})
	}
}
