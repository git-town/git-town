package configdomain

import (
	"slices"
	"strings"
)

// whether to display branch types in the CLI output
type DisplayTypes struct {
	Quantifier  Quantifier
	BranchTypes []BranchType
}

type Quantifier string

const (
	QuantifierAll  = "all"
	quantifierNo   = "no"
	quantifierOnly = ""
)

// indicates whether Git Town should display the given branch type
func (self DisplayTypes) ShouldDisplayType(branchType BranchType) bool {
	switch self.Quantifier {
	case QuantifierAll:
		return true
	case quantifierNo:
		return !slices.Contains(self.BranchTypes, branchType)
	case quantifierOnly:
		return slices.Contains(self.BranchTypes, branchType)
	}
	panic("unhandled DisplayType state: " + self.String())
}

func (self DisplayTypes) String() string {
	elements := []string{}
	if self.Quantifier != quantifierOnly {
		elements = append(elements, string(self.Quantifier))
	}
	for _, branchType := range self.BranchTypes {
		elements = append(elements, branchType.String())
	}
	return strings.Join(elements, " ")
}
