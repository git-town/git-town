package undoconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/undo/undoconfig"
	"github.com/git-town/git-town/v16/internal/undo/undodomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
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
			&opcodes.RemoveGlobalConfig{
				Key: configdomain.KeySyncPerennialStrategy,
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
				Added: []configdomain.Key(nil),
				Removed: map[configdomain.Key]string{
					configdomain.KeySyncPerennialStrategy: "1",
				},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []configdomain.Key(nil),
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.SetGlobalConfig{
				Key:   configdomain.KeySyncPerennialStrategy,
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
				Added:   []configdomain.Key(nil),
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []configdomain.Key(nil),
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.SetGlobalConfig{
				Key:   configdomain.KeyOffline,
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
			&opcodes.RemoveLocalConfig{
				Key: configdomain.KeySyncPerennialStrategy,
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
				Added:   []configdomain.Key(nil),
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added: []configdomain.Key(nil),
				Removed: map[configdomain.Key]string{
					configdomain.KeySyncPerennialStrategy: "1",
				},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.SetLocalConfig{
				Key:   configdomain.KeySyncPerennialStrategy,
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
				Added:   []configdomain.Key(nil),
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []configdomain.Key(nil),
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
			&opcodes.SetLocalConfig{
				Key:   configdomain.KeyOffline,
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
			&opcodes.RemoveGlobalConfig{
				Key: configdomain.KeySyncPerennialStrategy,
			},
			&opcodes.SetGlobalConfig{
				Key:   configdomain.KeyPushHook,
				Value: "0",
			},
			&opcodes.SetGlobalConfig{
				Key:   configdomain.KeyOffline,
				Value: "0",
			},
			&opcodes.RemoveLocalConfig{
				Key: configdomain.KeyPushHook,
			},
			&opcodes.SetLocalConfig{
				Key:   configdomain.KeyGithubToken,
				Value: "token",
			},
			&opcodes.SetLocalConfig{
				Key:   configdomain.KeyPerennialBranches,
				Value: "prod",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})
}

func emptyConfigDiff() undoconfig.ConfigDiff {
	return undoconfig.ConfigDiff{
		Added:   []configdomain.Key(nil),
		Changed: map[configdomain.Key]undodomain.Change[string]{},
		Removed: map[configdomain.Key]string{},
	}
}
