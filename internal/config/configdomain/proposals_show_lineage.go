package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// indicates whether and how proposals should display the stack lineage of the respective branch
type ProposalsShowLineage string

const (
	ProposalsShowLineageNone ProposalsShowLineage = "none" // don't display lineage in proposals
	ProposalsShowLineageCI   ProposalsShowLineage = "ci"   // lineage is embedded into proposals via https://github.com/git-town/action
	ProposalsShowLineageCLI  ProposalsShowLineage = "cli"  // the Git Town CLI should embed the lineage into proposals
)

func (self ProposalsShowLineage) IsCLI() bool {
	return self == ProposalsShowLineageCLI
}

func (self ProposalsShowLineage) String() string {
	return string(self)
}

func ParseProposalsShowLineage(value string) (Option[ProposalsShowLineage], error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		return None[ProposalsShowLineage](), nil
	case ProposalsShowLineageNone.String():
		return Some(ProposalsShowLineageNone), nil
	case ProposalsShowLineageCI.String():
		return Some(ProposalsShowLineageCI), nil
	case ProposalsShowLineageCLI.String():
		return Some(ProposalsShowLineageCLI), nil
	}
	parsedOpt, err := gohacks.ParseBoolOpt(value, "proposals-show-lineage")
	if err != nil {
		return None[ProposalsShowLineage](), fmt.Errorf(messages.ProposalsShowLineageInvalid, value)
	}
	if parsed, has := parsedOpt.Get(); has {
		if parsed {
			// The CLI is configured with "true" --> assume the user wants the CLI to embed lineage into proposals.
			return Some(ProposalsShowLineageCLI), nil
		}
		return Some(ProposalsShowLineageNone), nil
	}
	return None[ProposalsShowLineage](), fmt.Errorf(messages.ProposalsShowLineageInvalid, value)
}
