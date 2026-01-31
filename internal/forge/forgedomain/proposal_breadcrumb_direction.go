package forgedomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type ProposalBreadcrumbDirection string

const (
	ProposalBreadcrumbDirectionDown ProposalBreadcrumbDirection = "down"
	ProposalBreadcrumbDirectionUp   ProposalBreadcrumbDirection = "up"
)

func (self ProposalBreadcrumbDirection) String() string {
	return string(self)
}

func ParseProposalBreadcrumbDirection(value string, source string) (Option[ProposalBreadcrumbDirection], error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		return None[ProposalBreadcrumbDirection](), nil
	case ProposalBreadcrumbDirectionDown.String():
		return Some(ProposalBreadcrumbDirectionDown), nil
	case ProposalBreadcrumbDirectionUp.String():
		return Some(ProposalBreadcrumbDirectionUp), nil
	}
	return None[ProposalBreadcrumbDirection](), fmt.Errorf(messages.ProposalBreadcrumbDirectionInvalid, source, value)
}
