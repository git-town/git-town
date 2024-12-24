package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v17/internal/cli/colors"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

type BranchTypeOverrides map[gitdomain.LocalBranchName]BranchType

func (self BranchTypeOverrides) Add(other BranchTypeOverrides) BranchTypeOverrides {
	result := make(BranchTypeOverrides, len(self)+len(other))
	for key, value := range self {
		result[key] = value
	}
	for key, value := range other {
		result[key] = value
	}
	return result
}

func NewBranchTypeOverridesFromSnapshot(snapshot SingleSnapshot, removeLocalConfigValue removeLocalConfigValueFunc) (BranchTypeOverrides, error) {
	result := BranchTypeOverrides{}
	for key, value := range snapshot.BranchTypeOverrideEntries() {
		branch := key.Branch()
		if branch == "" {
			// empty branch --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigBranchTypeOverrideEmpty))
			_ = removeLocalConfigValue(key.Key)
			continue
		}
		value = strings.TrimSpace(value)
		if value == "" {
			// empty lineage entries are invalid --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigBranchTypeOverrideEmpty))
			_ = removeLocalConfigValue(key.Key)
			continue
		}
		branchTypeOpt, err := ParseBranchType(value)
		if err != nil {
			return result, err
		}
		if branchType, hasBranchType := branchTypeOpt.Get(); hasBranchType {
			result[branch] = branchType
		}
	}
	return result, nil
}
