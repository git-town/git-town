package undoconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/undo/undoconfig"
	"github.com/git-town/git-town/v24/internal/undo/undodomain"
	"github.com/git-town/git-town/v24/internal/vm/opcodes"
	"github.com/git-town/git-town/v24/internal/vm/program"
	"github.com/shoenig/test/must"
)

func TestConfigUndo(t *testing.T) {
	t.Parallel()

	t.Run("adding a value to the global cache", func(t *testing.T) {
		t.Parallel()
		before := configdomain.BeginConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
			Local:    configdomain.SingleSnapshot{},
			Unscoped: configdomain.SingleSnapshot{},
		}
		after := configdomain.EndConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline:               "0",
				configdomain.KeySyncPerennialStrategy: "1",
			},
			Local: configdomain.SingleSnapshot{},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeySyncPerennialStrategy,
				},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
				Removed: map[configdomain.Key]string{},
			},
			Local: emptyConfigDiff(),
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.ConfigRemove{
				Key:   configdomain.KeySyncPerennialStrategy,
				Scope: configdomain.ConfigScopeGlobal,
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("adding a value to the local cache", func(t *testing.T) {
		t.Parallel()
		before := configdomain.BeginConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
		}
		after := configdomain.EndConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyOffline:               "0",
				configdomain.KeySyncPerennialStrategy: "1",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: emptyConfigDiff(),
			Local: undoconfig.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeySyncPerennialStrategy,
				},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
				Removed: map[configdomain.Key]string{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.ConfigRemove{
				Key:   configdomain.KeySyncPerennialStrategy,
				Scope: configdomain.ConfigScopeLocal,
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("changing a value in the global cache", func(t *testing.T) {
		t.Parallel()
		before := configdomain.BeginConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
			Local: configdomain.SingleSnapshot{},
		}
		after := configdomain.EndConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "1",
			},
			Local: configdomain.SingleSnapshot{},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added: []configdomain.Key{},
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyOffline: {
						After:  "1",
						Before: "0",
					},
				},
				Removed: map[configdomain.Key]string{},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
				Removed: map[configdomain.Key]string{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.ConfigSet{
				Key:   configdomain.KeyOffline,
				Scope: configdomain.ConfigScopeGlobal,
				Value: "0",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("changing a value in the local cache", func(t *testing.T) {
		t.Parallel()
		before := configdomain.BeginConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
		}
		after := configdomain.EndConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "1",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
				Removed: map[configdomain.Key]string{},
			},
			Local: undoconfig.ConfigDiff{
				Added: []configdomain.Key{},
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyOffline: {
						After:  "1",
						Before: "0",
					},
				},
				Removed: map[configdomain.Key]string{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.ConfigSet{
				Key:   configdomain.KeyOffline,
				Scope: configdomain.ConfigScopeLocal,
				Value: "0",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("complex example", func(t *testing.T) {
		t.Parallel()
		before := configdomain.BeginConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline:  "0",
				configdomain.KeyPushHook: "0",
			},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyGithubToken:       "token",
				configdomain.KeyPerennialBranches: "prod",
			},
		}
		after := configdomain.EndConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline:               "1",
				configdomain.KeySyncPerennialStrategy: "1",
			},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyPerennialBranches: "prod qa",
				configdomain.KeyPushHook:          "1",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeySyncPerennialStrategy,
				},
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyOffline: {
						After:  "1",
						Before: "0",
					},
				},
				Removed: map[configdomain.Key]string{
					configdomain.KeyPushHook: "0",
				},
			},
			Local: undoconfig.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeyPushHook,
				},
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyPerennialBranches: {
						After:  "prod qa",
						Before: "prod",
					},
				},
				Removed: map[configdomain.Key]string{
					configdomain.KeyGithubToken: "token",
				},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.ConfigRemove{
				Key:   configdomain.KeySyncPerennialStrategy,
				Scope: configdomain.ConfigScopeGlobal,
			},
			&opcodes.ConfigSet{
				Key:   configdomain.KeyPushHook,
				Scope: configdomain.ConfigScopeGlobal,
				Value: "0",
			},
			&opcodes.ConfigSet{
				Key:   configdomain.KeyOffline,
				Scope: configdomain.ConfigScopeGlobal,
				Value: "0",
			},
			&opcodes.ConfigRemove{
				Key:   configdomain.KeyPushHook,
				Scope: configdomain.ConfigScopeLocal,
			},
			&opcodes.ConfigSet{
				Key:   configdomain.KeyGithubToken,
				Scope: configdomain.ConfigScopeLocal,
				Value: "token",
			},
			&opcodes.ConfigSet{
				Key:   configdomain.KeyPerennialBranches,
				Scope: configdomain.ConfigScopeLocal,
				Value: "prod",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("removing a value from the global cache", func(t *testing.T) {
		t.Parallel()
		before := configdomain.BeginConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline:               "0",
				configdomain.KeySyncPerennialStrategy: "1",
			},
			Local: configdomain.SingleSnapshot{},
		}
		after := configdomain.EndConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
			Local: configdomain.SingleSnapshot{},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
				Removed: map[configdomain.Key]string{
					configdomain.KeySyncPerennialStrategy: "1",
				},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
				Removed: map[configdomain.Key]string{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.ConfigSet{
				Key:   configdomain.KeySyncPerennialStrategy,
				Scope: configdomain.ConfigScopeGlobal,
				Value: "1",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("removing a value from the local cache", func(t *testing.T) {
		t.Parallel()
		before := configdomain.BeginConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyOffline:               "0",
				configdomain.KeySyncPerennialStrategy: "1",
			},
		}
		after := configdomain.EndConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
				Removed: map[configdomain.Key]string{},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
				Removed: map[configdomain.Key]string{
					configdomain.KeySyncPerennialStrategy: "1",
				},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.ConfigSet{
				Key:   configdomain.KeySyncPerennialStrategy,
				Scope: configdomain.ConfigScopeLocal,
				Value: "1",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})
}

func emptyConfigDiff() undoconfig.ConfigDiff {
	return undoconfig.ConfigDiff{
		Added:   []configdomain.Key{},
		Changed: map[configdomain.Key]undodomain.Change[string]{},
		Removed: map[configdomain.Key]string{},
	}
}
