package configdomain

import (
	"strings"

	"github.com/git-town/git-town/v23/internal/gohacks/slice"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/git-town/git-town/v23/pkg/set"
)

// ProposalBreadcrumbExcludeBranches lists the branch types to hide from proposal breadcrumbs.
type ProposalBreadcrumbExcludeBranches struct {
	set.Set[BranchType]
}

func (self ProposalBreadcrumbExcludeBranches) String() string {
	if len(self.Set) == 0 {
		return "(none)"
	}
	return strings.Join(slice.Stringify(self.Values()), ", ")
}

func NewProposalBreadcrumbExcludeBranches(branchTypes ...BranchType) ProposalBreadcrumbExcludeBranches {
	return ProposalBreadcrumbExcludeBranches{
		set.New(branchTypes...),
	}
}

func ParseProposalBreadcrumbExcludeBranches(text stringss.Trimmed, source string) (Option[ProposalBreadcrumbExcludeBranches], error) {
	parts := strings.Split(text.String(), ",")
	return ParseProposalBreadcrumbExcludeBranchesList(parts, source)
}

func ParseProposalBreadcrumbExcludeBranchesList(texts []string, source string) (Option[ProposalBreadcrumbExcludeBranches], error) {
	result := NewProposalBreadcrumbExcludeBranches()
	for _, text := range texts {
		branchTypeText := stringss.Trim(text)
		if branchTypeText == "" {
			continue
		}
		branchType, err := ParseBranchType(branchTypeText, source)
		if err != nil {
			return None[ProposalBreadcrumbExcludeBranches](), err
		}
		if branchTypeValue, hasBranchType := branchType.Get(); hasBranchType {
			result.Add(branchTypeValue)
		}
	}
	return Some(result), nil
}
