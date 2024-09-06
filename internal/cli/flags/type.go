package flags

import (
	"regexp"
	"strings"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	typeLong  = "type"
	typeShort = "t"
)

type branchTypeFlag Option[configdomain.BranchType]

// type-safe access to the CLI arguments of type configdomain.BranchType
func BranchType() (AddFunc, ReadTypeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringP(typeLong, typeShort, "", "limit the list of branches to switch to the given branch type(s)")
	}
	readFlag := func(cmd *cobra.Command) ([]configdomain.BranchType, error) {
		value, err := cmd.Flags().GetString(typeLong)
		if err != nil {
			panic(err)
		}
		return parseBranchTypes(value)
	}
	return addFlag, readFlag
}

func parseBranchTypes(text string) ([]configdomain.BranchType, error) {
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
	splitter := regexp.MustCompile(`,\+&|`)
	return splitter.Split(text, -1)
}

// the type signature for the function that reads the "type" flag from the args to the given Cobra command
type ReadTypeFlagFunc func(*cobra.Command) ([]configdomain.BranchType, error)

// // newEnum give a list of allowed flag parameters, where the second argument is the default
// func newBranchTypeFlag() branchTypeFlag {
// 	return branchTypeFlag(None[configdomain.BranchType]())
// }

// func (self branchTypeFlag) String() string {
// 	return self.String()
// }

// func (self *branchTypeFlag) Set(text string) error {
// 	branchType, validBranchType := configdomain.ParseBranchType(text).Get()
// 	if !validBranchType {
// 		return fmt.Errorf("invalid branch type: %q, allowed: contribution, feature, observed, parked, perennial, prototype", text)
// 	}
// 	*self = branchTypeFlag(Some(branchType))
// 	return nil
// }

// func (self *branchTypeFlag) Type() string {
// 	return "string"
// }
