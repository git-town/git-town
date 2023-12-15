package undo

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
)

// ConfigDiffs describes the changes made to the local and global Git configuration.
type ConfigDiffs struct {
	Global configdomain.ConfigDiff
	Local  configdomain.ConfigDiff
}

func NewConfigDiffs(before, after ConfigSnapshot) ConfigDiffs {
	globalCacheDiff := gitconfig.SingleCacheDiff(before.GitConfig.GlobalCache, after.GitConfig.GlobalCache)
	localCacheDiff := gitconfig.SingleCacheDiff(before.GitConfig.LocalCache, after.GitConfig.LocalCache)
	globalConfigDiff := configdomain.PartialConfigDiff(before.GitConfig.GlobalConfig, after.GitConfig.GlobalConfig)
	localConfigDiff := configdomain.PartialConfigDiff(before.GitConfig.LocalConfig, after.GitConfig.LocalConfig)
	fmt.Println("\n111111111111111111")
	spew.Dump(localCacheDiff)
	return ConfigDiffs{
		Global: globalCacheDiff.Merge(&globalConfigDiff),
		Local:  localCacheDiff.Merge(&localConfigDiff),
	}
}

func (self ConfigDiffs) UndoProgram() program.Program {
	result := program.Program{}
	for _, key := range self.Global.Added {
		result.Add(&opcode.RemoveGlobalConfig{Key: key})
	}
	for key, value := range self.Global.Removed {
		result.Add(&opcode.SetGlobalConfig{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range self.Global.Changed {
		result.Add(&opcode.SetGlobalConfig{
			Key:   key,
			Value: change.Before,
		})
	}
	for _, key := range self.Local.Added {
		result.Add(&opcode.RemoveLocalConfig{Key: key})
	}
	for key, value := range self.Local.Removed {
		result.Add(&opcode.SetLocalConfig{
			Key:   key,
			Value: value,
		})
	}
	for key, change := range self.Local.Changed {
		result.Add(&opcode.SetLocalConfig{
			Key:   key,
			Value: change.Before,
		})
	}
	return result
}
