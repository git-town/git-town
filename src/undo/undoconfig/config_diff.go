package undoconfig

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
)

// ConfigDiff describes changes made to the Git Town configuration.
type ConfigDiff struct {
	Added   []configdomain.Key
	Removed map[configdomain.Key]string
	Changed map[configdomain.Key]undodomain.Change[string]
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
	diff.Changed[key] = undodomain.Change[string]{
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
	diff.Changed[key] = undodomain.Change[string]{
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
		Changed: map[configdomain.Key]undodomain.Change[string]{},
	}
}
