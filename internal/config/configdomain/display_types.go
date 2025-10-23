package configdomain

import (
	"fmt"
	"slices"
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// whether to display branch types in the CLI output
type DisplayTypes struct {
	Quantifier  Quantifier
	BranchTypes []BranchType
}

type Quantifier string

const (
	QuantifierAll  = "all"
	QuantifierNo   = "no"
	QuantifierOnly = ""
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
	elements := []string{}
	if self.Quantifier != QuantifierOnly {
		elements = append(elements, string(self.Quantifier))
	}
	for _, branchType := range self.BranchTypes {
		elements = append(elements, branchType.String())
	}
	return strings.Join(elements, " ")
}

func ParseDisplayTypes(text, source string) (Option[DisplayTypes], error) {
	if len(text) == 0 {
		return None[DisplayTypes](), fmt.Errorf("please provide a valid entry for %s", source)
	}
	parts := strings.Split(text, " ")
	var quantifier Quantifier
	switch parts[0] {
	case QuantifierAll:
		quantifier = QuantifierAll
		parts = parts[1:]
		if len(parts) > 0 {
			return None[DisplayTypes](), fmt.Errorf(`the "all" quantifier for DisplayTypes does not accept branch types, you gave: %s`, parts)
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
		Quantifier:  quantifier,
		BranchTypes: branchTypes,
	}), nil
}
