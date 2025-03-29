package undoconfig

import (
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/gohacks/mapstools"
	"github.com/git-town/git-town/v18/internal/vm/opcodes"
	"github.com/git-town/git-town/v18/internal/vm/program"
)

// ConfigDiffs describes the changes made to the local and global Git configuration.
type ConfigDiffs struct {
	Global ConfigDiff
	Local  ConfigDiff
}

func NewConfigDiffs(before, after ConfigSnapshot) ConfigDiffs {
	return ConfigDiffs{
		Global: SingleCacheDiff(before.Global, after.Global),
		Local:  SingleCacheDiff(before.Local, after.Local),
	}
}

func (self ConfigDiffs) UndoProgram() program.Program {
	result := program.Program{}
	for _, key := range self.Global.Added {
		result.Add(&opcodes.ConfigRemove{
			Key:   key,
			Scope: configdomain.ConfigScopeGlobal,
		})
	}
	for _, removed := range mapstools.SortedKeyValues(self.Global.Removed) {
		result.Add(&opcodes.ConfigSet{
			Key:   removed.Key,
			Scope: configdomain.ConfigScopeGlobal,
			Value: removed.Value,
		})
	}
	for _, changed := range mapstools.SortedKeyValues(self.Global.Changed) {
		result.Add(&opcodes.ConfigSet{
			Key:   changed.Key,
			Scope: configdomain.ConfigScopeGlobal,
			Value: changed.Value.Before,
		})
	}
	for _, key := range self.Local.Added {
		result.Add(&opcodes.ConfigRemove{
			Key:   key,
			Scope: configdomain.ConfigScopeLocal,
		})
	}
	for _, removed := range mapstools.SortedKeyValues(self.Local.Removed) {
		result.Add(&opcodes.ConfigSet{
			Key:   removed.Key,
			Scope: configdomain.ConfigScopeLocal,
			Value: removed.Value,
		})
	}
	for _, changed := range mapstools.SortedKeyValues(self.Local.Changed) {
		result.Add(&opcodes.ConfigSet{
			Key:   changed.Key,
			Scope: configdomain.ConfigScopeLocal,
			Value: changed.Value.Before,
		})
	}
	return result
}
