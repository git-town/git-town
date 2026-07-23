package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/git-town/git-town/v24/internal/messages"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

type ProposalBreadcrumbDirection string

const (
	ProposalBreadcrumbDirectionDown ProposalBreadcrumbDirection = "down"
	ProposalBreadcrumbDirectionUp   ProposalBreadcrumbDirection = "up"
)

func (self ProposalBreadcrumbDirection) String() string {
	return string(self)
}

func ParseProposalBreadcrumbDirection(value stringss.Trimmed, source string) (Option[ProposalBreadcrumbDirection], error) {
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

func ParseProposalBreadcrumbDirectionOpt(valueOpt Option[string], source string) (Option[ProposalBreadcrumbDirection], error) {
	if value, has := valueOpt.Get(); has {
		return ParseProposalBreadcrumbDirection(value, source)
	}
	return None[ProposalBreadcrumbDirection](), nil
}
