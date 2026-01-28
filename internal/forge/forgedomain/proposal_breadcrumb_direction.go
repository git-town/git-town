package forgedomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type ProposalBreadcrumbDirection string

const (
	ProposalBreadcrumbDirectionTopDown  ProposalBreadcrumbDirection = "top-down"
	ProposalBreadcrumbDirectionBottomUp ProposalBreadcrumbDirection = "bottom-up"
)

func (self ProposalBreadcrumbDirection) String() string {
	return string(self)
}

func ParseProposalBreadcrumbDirection(value string, source string) (Option[ProposalBreadcrumbDirection], error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		return None[ProposalBreadcrumbDirection](), nil
	case ProposalBreadcrumbDirectionTopDown.String():
		return Some(ProposalBreadcrumbDirectionTopDown), nil
	case ProposalBreadcrumbDirectionBottomUp.String():
		return Some(ProposalBreadcrumbDirectionBottomUp), nil
	}
	return None[ProposalBreadcrumbDirection](), fmt.Errorf(messages.ProposalBreadcrumbDirectionInvalid, source, value)
}
