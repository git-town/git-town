package configdomain

import (
	"strings"

	"github.com/git-town/git-town/v23/internal/gohacks/slice"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/git-town/git-town/v23/pkg/set"
)

// ProposalBreadcrumbExclude lists the branch types to hide from proposal breadcrumbs.
type ProposalBreadcrumbExclude struct {
	set.Set[BranchType]
}

func (self ProposalBreadcrumbExclude) String() string {
	if len(self.Set) == 0 {
		return "(none)"
	}
	return strings.Join(slice.Stringify(self.Values()), ", ")
}

func NewProposalBreadcrumbExclude(branchTypes ...BranchType) ProposalBreadcrumbExclude {
	return ProposalBreadcrumbExclude{
		set.New(branchTypes...),
	}
}

func ParseProposalBreadcrumbExclude(text stringss.Trimmed, source string) (Option[ProposalBreadcrumbExclude], error) {
	parts := strings.Fields(text.String())
	return ParseProposalBreadcrumbExcludeList(parts, source)
}

func ParseProposalBreadcrumbExcludeList(texts []string, source string) (Option[ProposalBreadcrumbExclude], error) {
	result := NewProposalBreadcrumbExclude()
	for _, text := range texts {
		branchTypeText := stringss.Trim(text)
		if branchTypeText == "" {
			continue
		}
		branchType, err := ParseBranchType(branchTypeText, source)
		if err != nil {
			return None[ProposalBreadcrumbExclude](), err
		}
		if branchTypeValue, hasBranchType := branchType.Get(); hasBranchType {
			result.Add(branchTypeValue)
		}
	}
	return Some(result), nil
}
