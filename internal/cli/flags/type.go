package flags

import (
	"regexp"
	"sync"

	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/spf13/cobra"
)

const (
	typeLong  = "type"
	typeShort = "t"
)

// BranchType provides type-safe access to the CLI arguments of type configdomain.BranchType.
func BranchType() (AddFunc, ReadTypeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringP(typeLong, typeShort, "", "limit the list of branches to switch to the given branch type(s)")
	}
	readFlag := func(cmd *cobra.Command) ([]configdomain.BranchType, error) {
		value, err := cmd.Flags().GetString(typeLong)
		if err != nil {
			return []configdomain.BranchType{}, err
		}
		return ParseBranchTypes(stringss.Trim(value), "--type flag")
	}
	return addFlag, readFlag
}

func ParseBranchTypes(text stringss.Trimmed, source string) ([]configdomain.BranchType, error) {
	branchTypeNames := SplitBranchTypeNames(text)
	result := make([]configdomain.BranchType, 0, len(branchTypeNames))
	for _, branchTypeName := range branchTypeNames {
		branchTypeOpt, err := configdomain.ParseBranchType(branchTypeName, source)
		if err != nil {
			return result, err
		}
		if branchType, hasBranchType := branchTypeOpt.Get(); hasBranchType {
			result = append(result, branchType)
		}
	}
	return result, nil
}

func SplitBranchTypeNames(text stringss.Trimmed) []stringss.Trimmed {
	splitBranchOnce.Do(func() {
		splitBranchRegex = regexp.MustCompile(`[,\+&\|]`)
	})
	splitted := splitBranchRegex.Split(text.String(), -1)
	result := make([]stringss.Trimmed, 0, len(splitted))
	for _, split := range splitted {
		trimmedSplit := stringss.Trim(split)
		if len(trimmedSplit) > 0 {
			result = append(result, trimmedSplit)
		}
	}
	return result
}

var (
	splitBranchOnce  sync.Once
	splitBranchRegex *regexp.Regexp
)

// ReadTypeFlagFunc is the type signature for the function that reads the "type" flag from the args to the given Cobra command.
type ReadTypeFlagFunc func(*cobra.Command) ([]configdomain.BranchType, error)
