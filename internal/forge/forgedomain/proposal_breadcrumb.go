package forgedomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ProposalBreadcrumb indicates whether and how proposals should display the stack lineage of the respective branch.
type ProposalBreadcrumb string

const (
	ProposalBreadcrumbNone     ProposalBreadcrumb = "none"     // don't display lineage in proposals
	ProposalBreadcrumbStacks   ProposalBreadcrumb = "stacks"   // the Git Town CLI should embed breadcrumbs into proposals for stacks with more than one branch
	ProposalBreadcrumbBranches ProposalBreadcrumb = "branches" // the Git Town CLI should embed breadcrumbs into proposals for all branches
)

// EmbedBreadcrumb indicates whether the Git Town CLI should embed the breadcrumb into proposals.
func (self ProposalBreadcrumb) DisplayBreadcrumb(stackSize int) bool {
	switch self {
	case ProposalBreadcrumbBranches:
		return true
	case ProposalBreadcrumbNone:
		return false
	case ProposalBreadcrumbStacks:
		return stackSize > 2
	}
	panic(fmt.Sprintf("unexpected ProposalBreadcrumb: %#v", self))
}

func (self ProposalBreadcrumb) Enabled() bool {
	switch self {
	case ProposalBreadcrumbBranches, ProposalBreadcrumbStacks:
		return true
	case ProposalBreadcrumbNone:
		return false
	}
	panic(fmt.Sprintf("unexpected ProposalBreadcrumb: %#v", self))
}

func (self ProposalBreadcrumb) String() string {
	return string(self)
}

func ParseProposalBreadcrumb(value string, source string) (Option[ProposalBreadcrumb], error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		return None[ProposalBreadcrumb](), nil
	case ProposalBreadcrumbNone.String():
		return Some(ProposalBreadcrumbNone), nil
	case ProposalBreadcrumbBranches.String():
		return Some(ProposalBreadcrumbBranches), nil
	}
	parsedOpt, err := gohacks.ParseBoolOpt[bool](value, "proposal-breadcrumb")
	if err != nil {
		return None[ProposalBreadcrumb](), fmt.Errorf(messages.ProposalBreadcrumbInvalid, source, value)
	}
	if parsed, has := parsedOpt.Get(); has {
		if parsed {
			// The CLI is configured with "true" --> assume the user wants the CLI to embed lineage into proposals.
			return Some(ProposalBreadcrumbBranches), nil
		}
		return Some(ProposalBreadcrumbNone), nil
	}
	return None[ProposalBreadcrumb](), fmt.Errorf(messages.ProposalBreadcrumbInvalid, source, value)
}
