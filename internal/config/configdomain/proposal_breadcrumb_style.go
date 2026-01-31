package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type ProposalBreadcrumbStyle string

const (
	ProposalBreadcrumbStyleTree ProposalBreadcrumbStyle = "tree" // always render the breadcrumb as a tree, even if linear
	ProposalBreadcrumbStyleAuto ProposalBreadcrumbStyle = "auto" // render the breadcrumb flat if linear, otherwise as a tree
)

func (self ProposalBreadcrumbStyle) String() string {
	return string(self)
}

func ParseProposalBreadcrumbStyle(value string, source string) (Option[ProposalBreadcrumbStyle], error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		return None[ProposalBreadcrumbStyle](), nil
	case ProposalBreadcrumbStyleTree.String():
		return Some(ProposalBreadcrumbStyleTree), nil
	case ProposalBreadcrumbStyleAuto.String():
		return Some(ProposalBreadcrumbStyleAuto), nil
	}
	return None[ProposalBreadcrumbStyle](), fmt.Errorf(messages.ProposalBreadcrumbStyleInvalid, source, value)
}
