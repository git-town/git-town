package configdomain

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// whether to display branch types in the CLI output
type DisplayTypes struct {
	BranchTypes []BranchType
	Quantifier  Quantifier
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
	switch self.Quantifier {
	case QuantifierAll:
		return "all branch types"
	case QuantifierNo:
		if len(self.BranchTypes) == 0 {
			return "no branch types"
		}
		return "all branch types except " + formatBranchTypeList(self.BranchTypes)
	case QuantifierOnly:
		return "only the branch types " + formatBranchTypeList(self.BranchTypes)
	}
	panic("unhandled DisplayType state: " + self.Quantifier)
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

func formatBranchTypeList(branchTypes []BranchType) string {
	if len(branchTypes) == 0 {
		return ""
	}
	if len(branchTypes) == 1 {
		return `"` + branchTypes[0].String() + `"`
	}
	if len(branchTypes) == 2 {
		return `"` + branchTypes[0].String() + `" and "` + branchTypes[1].String() + `"`
	}
	// 3+ items: use Oxford comma style
	var parts []string
	for i, branchType := range branchTypes {
		if i == len(branchTypes)-1 {
			parts = append(parts, `and "`+branchType.String()+`"`)
		} else {
			parts = append(parts, `"`+branchType.String()+`"`)
		}
	}
	return strings.Join(parts, ", ")
}
