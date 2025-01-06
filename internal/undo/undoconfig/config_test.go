package undoconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/undo/undoconfig"
	"github.com/git-town/git-town/v17/internal/undo/undodomain"
	"github.com/git-town/git-town/v17/internal/vm/opcodes"
	"github.com/git-town/git-town/v17/internal/vm/program"
	"github.com/shoenig/test/must"
)

func TestConfigUndo(t *testing.T) {
	t.Parallel()

	t.Run("adding a value to the global cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
			Local: configdomain.SingleSnapshot{},
		}
		after := undoconfig.ConfigSnapshot{
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
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
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

	t.Run("removing a value from the global cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline:               "0",
				configdomain.KeySyncPerennialStrategy: "1",
			},
			Local: configdomain.SingleSnapshot{},
		}
		after := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
			Local: configdomain.SingleSnapshot{},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added: []configdomain.Key{},
				Removed: map[configdomain.Key]string{
					configdomain.KeySyncPerennialStrategy: "1",
				},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
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

	t.Run("changing a value in the global cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
			Local: configdomain.SingleSnapshot{},
		}
		after := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "1",
			},
			Local: configdomain.SingleSnapshot{},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
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

	t.Run("adding a value to the local cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
		}
		after := undoconfig.ConfigSnapshot{
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
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
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

	t.Run("removing a value from the local cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyOffline:               "0",
				configdomain.KeySyncPerennialStrategy: "1",
			},
		}
		after := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added: []configdomain.Key{},
				Removed: map[configdomain.Key]string{
					configdomain.KeySyncPerennialStrategy: "1",
				},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
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

	t.Run("changing a value in the local cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "0",
			},
		}
		after := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyOffline: "1",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
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
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				configdomain.KeyOffline:  "0",
				configdomain.KeyPushHook: "0",
			},
			Local: configdomain.SingleSnapshot{
				configdomain.KeyPerennialBranches: "prod",
				configdomain.KeyGithubToken:       "token",
			},
		}
		after := undoconfig.ConfigSnapshot{
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
				Removed: map[configdomain.Key]string{
					configdomain.KeyPushHook: "0",
				},
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: undoconfig.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeyPushHook,
				},
				Removed: map[configdomain.Key]string{
					configdomain.KeyGithubToken: "token",
				},
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyPerennialBranches: {
						Before: "prod",
						After:  "prod qa",
					},
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
}

func emptyConfigDiff() undoconfig.ConfigDiff {
	return undoconfig.ConfigDiff{
		Added:   []configdomain.Key{},
		Changed: map[configdomain.Key]undodomain.Change[string]{},
		Removed: map[configdomain.Key]string{},
	}
}
