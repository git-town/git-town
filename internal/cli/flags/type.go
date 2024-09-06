package flags

import (
	"fmt"

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
		cmd.PersistentFlags().StringP(typeLong, typeShort, "", "display all Git commands run under the hood")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.BranchType], error) {
		value, err := cmd.Flags().GetString(typeLong)
		if err != nil {
			panic(err)
		}
		branchType, err := parseBranchType(value)
		if err != nil {
			return None[configdomain.BranchType](), err
		}
		return branchType, nil
	}
	return addFlag, readFlag
}

func parseBranchType(text string) (branchType Option[configdomain.BranchType], err error) {
	switch text {
	case "contribution":
		return Some(configdomain.BranchTypeContributionBranch), nil
	case "feature":
		return Some(configdomain.BranchTypeFeatureBranch), nil
	case "observed":
		return Some(configdomain.BranchTypeObservedBranch), nil
	case "parked":
		return Some(configdomain.BranchTypeParkedBranch), nil
	case "perennial":
		return Some(configdomain.BranchTypePerennialBranch), nil
	case "prototype":
		return Some(configdomain.BranchTypePrototypeBranch), nil
	case "":
		return None[configdomain.BranchType](), nil
	}
	return None[configdomain.BranchType](), fmt.Errorf("invalid branch type: %q, allowed: contribution, feature, observed, parked, perennial, prototype", text)
}

// the type signature for the function that reads the "type" flag from the args to the given Cobra command
type ReadTypeFlagFunc func(*cobra.Command) (Option[configdomain.BranchType], error)

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
