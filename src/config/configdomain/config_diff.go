package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/google/go-cmp/cmp"
)

// ConfigDiff describes changes made to the Git Town configuration.
type ConfigDiff struct {
	Added   []Key
	Removed map[Key]string
	Changed map[Key]domain.Change[string]
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
func DiffLocalBranchNames(diff *ConfigDiff, key Key, beforeValue *domain.LocalBranchNames, afterValue *domain.LocalBranchNames) {
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
	diff.Changed[key] = domain.Change[string]{
		Before: beforeValue.String(),
		After:  afterValue.String(),
	}
}

type diffArg interface {
	fmt.Stringer
	comparable
}

func DiffPtr[T diffArg](diff *ConfigDiff, key Key, before *T, after *T) {
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
