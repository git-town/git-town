package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/google/go-cmp/cmp"
)

// ConfigDiff describes the changes made to the Git Town configuration.
type ConfigDiff struct {
	Added   []Key
	Removed map[Key]string
	Changed map[Key]domain.Change[string]
}

// Merge merges the given ConfigDiff into this ConfigDiff.
func (self *ConfigDiff) Merge(other *ConfigDiff) ConfigDiff {
	result := EmptyConfigDiff()
	result.Added = append(self.Added, other.Added...) //nolint:gocritic
	for key, value := range self.Removed {
		result.Removed[key] = value
	}
	for key, value := range other.Removed {
		result.Removed[key] = value
	}
	for key, value := range self.Changed {
		result.Changed[key] = value
	}
	for key, value := range other.Changed {
		result.Changed[key] = value
	}
	return result
}

type checkArg interface {
	fmt.Stringer
	comparable
}

func Check[T checkArg](diff *ConfigDiff, key Key, before T, after T) {
	beforeText := before.String()
	afterText := after.String()
	if beforeText == afterText {
		return
	}
	if beforeText == "" {
		diff.Added = append(diff.Added, key)
		return
	}
	if afterText == "" {
		diff.Removed[key] = beforeText
		return
	}
	diff.Changed[key] = domain.Change[string]{
		Before: beforeText,
		After:  afterText,
	}
}

func DiffLocalBranchNames(diff *ConfigDiff, key Key, before *domain.LocalBranchNames, after *domain.LocalBranchNames) {
	if cmp.Equal(before, after) {
		return
	}
	if before == nil || len(*before) == 0 {
		diff.Added = append(diff.Added, key)
		return
	}
	if after == nil || len(*after) == 0 {
		diff.Removed[key] = before.String()
	}
	diff.Changed[key] = domain.Change[string]{
		Before: before.String(),
		After:  after.String(),
	}
}

func DiffPtr[T checkArg](diff *ConfigDiff, key Key, before *T, after *T) {
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
	diff.Changed[key] = domain.Change[string]{
		Before: (*before).String(),
		After:  (*after).String(),
	}
}

func DiffString(diff *ConfigDiff, key Key, before string, after string) {
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
	diff.Changed[key] = domain.Change[string]{
		Before: before,
		After:  after,
	}
}

func DiffStringPtr(diff *ConfigDiff, key Key, before *string, after *string) {
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
		Added:   []Key{},
		Removed: map[Key]string{},
		Changed: map[Key]domain.Change[string]{},
	}
}
