package undodomain

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/google/go-cmp/cmp"
)

// ConfigDiff describes changes made to the Git Town configuration.
type ConfigDiff struct {
	Added   []configdomain.Key
	Removed map[configdomain.Key]string
	Changed map[configdomain.Key]Change[string]
}

// Merge merges the given ConfigDiff into this ConfigDiff, overwriting values that exist in both.
func (self *ConfigDiff) Merge(other *ConfigDiff) {
	self.Added = append(self.Added, other.Added...)
	for key, value := range other.Removed {
		self.Removed[key] = value
	}
	for key, value := range other.Changed {
		self.Changed[key] = value
	}
}

// DiffLocalBranchNames adds the difference between the given before and after values
// for the attribute with the given key to the given ConfigDiff.
func DiffLocalBranchNames(diff *ConfigDiff, key configdomain.Key, beforeValue *gitdomain.LocalBranchNames, afterValue *gitdomain.LocalBranchNames) {
	if cmp.Equal(beforeValue, afterValue) {
		return
	}
	if beforeValue == nil || len(*beforeValue) == 0 {
		diff.Added = append(diff.Added, key)
		return
	}
	if afterValue == nil || len(*afterValue) == 0 {
		diff.Removed[key] = beforeValue.String()
	}
	diff.Changed[key] = Change[string]{
		Before: beforeValue.String(),
		After:  afterValue.String(),
	}
}

type diffArg interface {
	fmt.Stringer
	comparable
}

func DiffPtr[T diffArg](diff *ConfigDiff, key configdomain.Key, before *T, after *T) {
	if before == nil && after == nil {
		return
	}
	if before == nil {
		diff.Added = append(diff.Added, key)
		return
	}
	if after == nil {
		diff.Removed[key] = (*before).String()
		return
	}
	if *before == *after {
		return
	}
	diff.Changed[key] = Change[string]{
		Before: (*before).String(),
		After:  (*after).String(),
	}
}

func DiffString(diff *ConfigDiff, key configdomain.Key, before string, after string) {
	if before == after {
		return
	}
	if before == "" {
		diff.Added = append(diff.Added, key)
		return
	}
	if after == "" {
		diff.Removed[key] = before
		return
	}
	diff.Changed[key] = Change[string]{
		Before: before,
		After:  after,
	}
}

func DiffStringPtr(diff *ConfigDiff, key configdomain.Key, before *string, after *string) {
	beforeText := ""
	if before != nil {
		beforeText = *before
	}
	afterText := ""
	if after != nil {
		afterText = *after
	}
	DiffString(diff, key, beforeText, afterText)
}

func EmptyConfigDiff() ConfigDiff {
	return ConfigDiff{
		Added:   []configdomain.Key{},
		Removed: map[configdomain.Key]string{},
		Changed: map[configdomain.Key]Change[string]{},
	}
}

// PartialConfigDiff diffs the given PartialConfig instances.
func PartialConfigDiff(before, after configdomain.PartialConfig) ConfigDiff {
	result := ConfigDiff{
		Added:   []configdomain.Key{},
		Removed: map[configdomain.Key]string{},
		Changed: map[configdomain.Key]Change[string]{},
	}
	DiffPtr(&result, configdomain.KeyCodeHostingOriginHostname, before.CodeHostingOriginHostname, after.CodeHostingOriginHostname)
	DiffPtr(&result, configdomain.KeyCodeHostingPlatform, before.CodeHostingPlatformName, after.CodeHostingPlatformName)
	DiffPtr(&result, configdomain.KeyGiteaToken, before.GiteaToken, after.GiteaToken)
	DiffPtr(&result, configdomain.KeyGithubToken, before.GitHubToken, after.GitHubToken)
	DiffPtr(&result, configdomain.KeyGitlabToken, before.GitLabToken, after.GitLabToken)
	DiffPtr(&result, configdomain.KeyMainBranch, before.MainBranch, after.MainBranch)
	DiffPtr(&result, configdomain.KeyOffline, before.Offline, after.Offline)
	DiffLocalBranchNames(&result, configdomain.KeyPerennialBranches, before.PerennialBranches, after.PerennialBranches)
	DiffPtr(&result, configdomain.KeyPushHook, before.PushHook, after.PushHook)
	DiffPtr(&result, configdomain.KeyPushNewBranches, before.NewBranchPush, after.NewBranchPush)
	DiffPtr(&result, configdomain.KeyShipDeleteTrackingBranch, before.ShipDeleteTrackingBranch, after.ShipDeleteTrackingBranch)
	DiffPtr(&result, configdomain.KeySyncFeatureStrategy, before.SyncFeatureStrategy, after.SyncFeatureStrategy)
	DiffPtr(&result, configdomain.KeySyncPerennialStrategy, before.SyncPerennialStrategy, after.SyncPerennialStrategy)
	DiffPtr(&result, configdomain.KeySyncUpstream, before.SyncUpstream, after.SyncUpstream)
	return result
}
