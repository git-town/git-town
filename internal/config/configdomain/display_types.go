package configdomain

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// whether to display branch types in the CLI output
type DisplayTypes struct {
	BranchTypes []BranchType // the branch types for which the user has specified exceptions
	Quantifier  Quantifier   // whether to include or exclude the listed branches
}

// Quantifier specifies whether to display or not display branches
type Quantifier string

const (
	QuantifierAll  = "all" // display all branches
	QuantifierNo   = "no"  // display no or all except the specified branches
	QuantifierOnly = ""    // display only the specified branches
)

// indicates whether Git Town should display the given branch type
func (self DisplayTypes) ShouldDisplayType(branchType BranchType) bool {
	switch self.Quantifier {
	case QuantifierAll:
		return true
	case QuantifierNo:
		if len(self.BranchTypes) == 0 {
			return false
		}
		return !slices.Contains(self.BranchTypes, branchType)
	case QuantifierOnly:
		return slices.Contains(self.BranchTypes, branchType)
	}
	panic("unhandled DisplayType state: " + self.String())
}

func (self DisplayTypes) String() string {
	switch self.Quantifier {
	case QuantifierAll:
		return "all branch types"
	case QuantifierNo:
		if len(self.BranchTypes) == 0 {
			return "no branch types"
		}
		return "all branch types except " + slice.JoinSentenceQuotes(self.BranchTypes)
	case QuantifierOnly:
		return "only the branch types " + slice.JoinSentenceQuotes(self.BranchTypes)
	}
	panic("unhandled DisplayType quantifier: " + self.Quantifier)
}

func ParseDisplayTypes(text, source string) (Option[DisplayTypes], error) {
	if len(text) == 0 {
		return None[DisplayTypes](), nil
	}
	re := regexp.MustCompile(`[ +\-&_]`)
	parts := re.Split(strings.ToLower(text), -1)
	var quantifier Quantifier
	switch parts[0] {
	case QuantifierAll:
		quantifier = QuantifierAll
		parts = parts[1:]
		if len(parts) > 0 {
			return None[DisplayTypes](), fmt.Errorf(`the "all" quantifier for DisplayTypes does not accept branch types, in %q you gave: %s`, source, parts)
		}
	case QuantifierNo:
		quantifier = QuantifierNo
		parts = parts[1:]
	default:
		quantifier = QuantifierOnly
	}
	branchTypes := make([]BranchType, len(parts))
	for p, part := range parts {
		branchTypeOpt, err := ParseBranchType(part)
		if err != nil {
			return None[DisplayTypes](), err
		}
		if branchType, hasBranchType := branchTypeOpt.Get(); hasBranchType {
			branchTypes[p] = branchType
		}
	}
	return Some(DisplayTypes{
		BranchTypes: branchTypes,
		Quantifier:  quantifier,
	}), nil
}
