package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v17/internal/cli/colors"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

type BranchTypeOverrides map[gitdomain.LocalBranchName]BranchType

func NewBranchTypeOverridesFromSnapshot(snapshot SingleSnapshot) (BranchTypeOverrides, error) {
	result := BranchTypeOverrides{}
	for key, value := range snapshot.BranchTypeOverrideEntries() {
		childName := key.ChildName()
		if childName == "" {
			// empty lineage entries are invalid --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigLineageEmptyChild))
			_ = removeLocalConfigValue(key.Key())
			continue
		}
		child := gitdomain.NewLocalBranchName(childName)
		value = strings.TrimSpace(value)
		if value == "" {
			// empty lineage entries are invalid --> delete it
			fmt.Println(colors.Cyan().Styled(messages.ConfigLineageEmptyChild))
			_ = removeLocalConfigValue(key.Key())
			continue
		}
		if updateOutdated && childName == value {
			fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.ConfigLineageParentIsChild, childName)))
			_ = removeLocalConfigValue(NewParentKey(child))
		}
		parent := gitdomain.NewLocalBranchName(value)
		result = result.Set(child, parent)
	}
	return result, nil
}
