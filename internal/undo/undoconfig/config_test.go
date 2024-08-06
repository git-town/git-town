package undoconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/internal/undo/undoconfig"
	"github.com/git-town/git-town/v14/internal/undo/undodomain"
	"github.com/git-town/git-town/v14/internal/vm/opcodes"
	"github.com/git-town/git-town/v14/internal/vm/program"
	"github.com/git-town/git-town/v14/pkg/keys"
	"github.com/shoenig/test/must"
)

func TestConfigUndo(t *testing.T) {
	t.Parallel()

	t.Run("adding a value to the global cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				keys.KeyOffline: "0",
			},
			Local: configdomain.SingleSnapshot{},
		}
		after := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				keys.KeyOffline:               "0",
				keys.KeySyncPerennialStrategy: "1",
			},
			Local: configdomain.SingleSnapshot{},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added: []keys.Key{
					keys.KeySyncPerennialStrategy,
				},
				Removed: map[keys.Key]string{},
				Changed: map[keys.Key]undodomain.Change[string]{},
			},
			Local: emptyConfigDiff(),
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.RemoveGlobalConfig{
				Key: keys.KeySyncPerennialStrategy,
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("removing a value from the global cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				keys.KeyOffline:               "0",
				keys.KeySyncPerennialStrategy: "1",
			},
			Local: configdomain.SingleSnapshot{},
		}
		after := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				keys.KeyOffline: "0",
			},
			Local: configdomain.SingleSnapshot{},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added: []keys.Key{},
				Removed: map[keys.Key]string{
					keys.KeySyncPerennialStrategy: "1",
				},
				Changed: map[keys.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []keys.Key{},
				Removed: map[keys.Key]string{},
				Changed: map[keys.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.SetGlobalConfig{
				Key:   keys.KeySyncPerennialStrategy,
				Value: "1",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("changing a value in the global cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				keys.KeyOffline: "0",
			},
			Local: configdomain.SingleSnapshot{},
		}
		after := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				keys.KeyOffline: "1",
			},
			Local: configdomain.SingleSnapshot{},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []keys.Key{},
				Removed: map[keys.Key]string{},
				Changed: map[keys.Key]undodomain.Change[string]{
					keys.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []keys.Key{},
				Removed: map[keys.Key]string{},
				Changed: map[keys.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.SetGlobalConfig{
				Key:   keys.KeyOffline,
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
				keys.KeyOffline: "0",
			},
		}
		after := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				keys.KeyOffline:               "0",
				keys.KeySyncPerennialStrategy: "1",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: emptyConfigDiff(),
			Local: undoconfig.ConfigDiff{
				Added: []keys.Key{
					keys.KeySyncPerennialStrategy,
				},
				Removed: map[keys.Key]string{},
				Changed: map[keys.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.RemoveLocalConfig{
				Key: keys.KeySyncPerennialStrategy,
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("removing a value from the local cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				keys.KeyOffline:               "0",
				keys.KeySyncPerennialStrategy: "1",
			},
		}
		after := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				keys.KeyOffline: "0",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []keys.Key{},
				Removed: map[keys.Key]string{},
				Changed: map[keys.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added: []keys.Key{},
				Removed: map[keys.Key]string{
					keys.KeySyncPerennialStrategy: "1",
				},
				Changed: map[keys.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.SetLocalConfig{
				Key:   keys.KeySyncPerennialStrategy,
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
				keys.KeyOffline: "0",
			},
		}
		after := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{},
			Local: configdomain.SingleSnapshot{
				keys.KeyOffline: "1",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []keys.Key{},
				Removed: map[keys.Key]string{},
				Changed: map[keys.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []keys.Key{},
				Removed: map[keys.Key]string{},
				Changed: map[keys.Key]undodomain.Change[string]{
					keys.KeyOffline: {
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
				Key:   keys.KeyOffline,
				Value: "0",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("complex example", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				keys.KeyOffline:  "0",
				keys.KeyPushHook: "0",
			},
			Local: configdomain.SingleSnapshot{
				keys.KeyPerennialBranches: "prod",
				keys.KeyGithubToken:       "token",
			},
		}
		after := undoconfig.ConfigSnapshot{
			Global: configdomain.SingleSnapshot{
				keys.KeyOffline:               "1",
				keys.KeySyncPerennialStrategy: "1",
			},
			Local: configdomain.SingleSnapshot{
				keys.KeyPerennialBranches: "prod qa",
				keys.KeyPushHook:          "1",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added: []keys.Key{
					keys.KeySyncPerennialStrategy,
				},
				Removed: map[keys.Key]string{
					keys.KeyPushHook: "0",
				},
				Changed: map[keys.Key]undodomain.Change[string]{
					keys.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: undoconfig.ConfigDiff{
				Added: []keys.Key{
					keys.KeyPushHook,
				},
				Removed: map[keys.Key]string{
					keys.KeyGithubToken: "token",
				},
				Changed: map[keys.Key]undodomain.Change[string]{
					keys.KeyPerennialBranches: {
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
				Key: keys.KeySyncPerennialStrategy,
			},
			&opcodes.SetGlobalConfig{
				Key:   keys.KeyPushHook,
				Value: "0",
			},
			&opcodes.SetGlobalConfig{
				Key:   keys.KeyOffline,
				Value: "0",
			},
			&opcodes.RemoveLocalConfig{
				Key: keys.KeyPushHook,
			},
			&opcodes.SetLocalConfig{
				Key:   keys.KeyGithubToken,
				Value: "token",
			},
			&opcodes.SetLocalConfig{
				Key:   keys.KeyPerennialBranches,
				Value: "prod",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})
}

func emptyConfigDiff() undoconfig.ConfigDiff {
	return undoconfig.ConfigDiff{
		Added:   []keys.Key{},
		Changed: map[keys.Key]undodomain.Change[string]{},
		Removed: map[keys.Key]string{},
	}
}
