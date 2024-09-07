package flags

import (
	"regexp"
	"strings"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const (
	typeLong  = "type"
	typeShort = "t"
)

// type-safe access to the CLI arguments of type configdomain.BranchType
func BranchType() (AddFunc, ReadTypeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringP(typeLong, typeShort, "", "limit the list of branches to switch to the given branch type(s)")
	}
	readFlag := func(cmd *cobra.Command) ([]configdomain.BranchType, error) {
		value, err := cmd.Flags().GetString(typeLong)
		if err != nil {
			return []configdomain.BranchType{}, err
		}
		return ParseBranchTypes(value)
	}
	return addFlag, readFlag
}

func ParseBranchTypes(text string) ([]configdomain.BranchType, error) {
	branchTypeNames := SplitBranchTypeNames(text)
	result := make([]configdomain.BranchType, 0, len(branchTypeNames))
	for _, branchTypeName := range branchTypeNames {
		branchTypeOpt, err := configdomain.ParseBranchType(branchTypeName)
		if err != nil {
			return result, err
		}
		if branchType, hasBranchType := branchTypeOpt.Get(); hasBranchType {
			result = append(result, branchType)
		}
	}
	return result, nil
}

func SplitBranchTypeNames(text string) []string {
	text = strings.TrimSpace(text)
	splitter := regexp.MustCompile(`[,\+&\|]`)
	splitted := splitter.Split(text, -1)
	result := make([]string, 0, len(splitted))
	for _, split := range splitted {
		if len(split) > 0 {
			result = append(result, split)
		}
	}
	return result
}

// the type signature for the function that reads the "type" flag from the args to the given Cobra command
type ReadTypeFlagFunc func(*cobra.Command) ([]configdomain.BranchType, error)
