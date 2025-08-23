package format

import (
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

const indent = "  "

// BranchLineage provides printable formatting of the given branch lineage.
func BranchLineage(lineage configdomain.Lineage) string {
	result := strings.Builder{}
	for _, root := range lineage.Roots() {
		branchTree(branchTreeArgs{
			branch:  root,
			indent:  0,
			lineage: lineage,
			builder: &result,
		})
	}
	return result.String()
}

type branchTreeArgs struct {
	branch  gitdomain.LocalBranchName
	indent  int
	lineage configdomain.Lineage
	builder *strings.Builder
}

// branchTree provids a printable version of the given branch tree.
func branchTree(args branchTreeArgs) {
	for range args.indent {
		args.builder.WriteString(indent)
	}
	args.builder.WriteString(args.branch.String())
	for _, child := range args.lineage.Children(args.branch) {
		args.builder.WriteString("\n")
		branchTree(branchTreeArgs{
			branch:  child,
			indent:  args.indent + 1,
			lineage: args.lineage,
			builder: args.builder,
		})
	}
}
