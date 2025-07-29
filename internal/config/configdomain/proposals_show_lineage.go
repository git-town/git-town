package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// indicates whether and how proposals should display the stack lineage of the respective branch
type ProposalsShowLineage string

const (
	ProposalsShowLineageNone ProposalsShowLineage = "none" // don't display lineage in proposals
	ProposalsShowLineageCI   ProposalsShowLineage = "ci"   // this team has set up https://github.com/git-town/action to embed the stack lineage into proposals
	ProposalsShowLineageCLI  ProposalsShowLineage = "cli"  // the Git Town CLI should embed the lineage into proposals
)

func (self ProposalsShowLineage) String() string {
	return string(self)
}

func ParseProposalsShowLineage(value string) (Option[ProposalsShowLineage], error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		return None[ProposalsShowLineage](), nil
	case "ci":
		return Some(ProposalsShowLineageCI), nil
	case "cli":
		return Some(ProposalsShowLineageCLI), nil
	default:
		return None[ProposalsShowLineage](), fmt.Errorf(messages.ProposalsShowLineageInvalid, value)
	}
}
