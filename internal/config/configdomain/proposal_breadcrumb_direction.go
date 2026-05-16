package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type ProposalBreadcrumbDirection string

const (
	ProposalBreadcrumbDirectionDown ProposalBreadcrumbDirection = "down"
	ProposalBreadcrumbDirectionUp   ProposalBreadcrumbDirection = "up"
)

func (self ProposalBreadcrumbDirection) String() string {
	return string(self)
}

func ParseProposalBreadcrumbDirection(value stringss.TrimmedString, source string) (Option[ProposalBreadcrumbDirection], error) {
	switch strings.ToLower(value.String()) {
	case "":
		return None[ProposalBreadcrumbDirection](), nil
	case ProposalBreadcrumbDirectionDown.String():
		return Some(ProposalBreadcrumbDirectionDown), nil
	case ProposalBreadcrumbDirectionUp.String():
		return Some(ProposalBreadcrumbDirectionUp), nil
	}
	return None[ProposalBreadcrumbDirection](), fmt.Errorf(messages.ProposalBreadcrumbDirectionInvalid, source, value)
}
