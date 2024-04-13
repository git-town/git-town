package undoconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/undo/undodomain"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/shoenig/test/must"
)

func TestConfigUndo(t *testing.T) {
	t.Parallel()

	t.Run("adding a value to the global cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline: "0",
			},
			Local: gitconfig.SingleSnapshot{},
		}
		after := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline:               "0",
				gitconfig.KeySyncPerennialStrategy: "1",
			},
			Local: gitconfig.SingleSnapshot{},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added: []gitconfig.Key{
					gitconfig.KeySyncPerennialStrategy,
				},
				Removed: map[gitconfig.Key]string{},
				Changed: map[gitconfig.Key]undodomain.Change[string]{},
			},
			Local: emptyConfigDiff(),
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.RemoveGlobalConfig{
				Key: gitconfig.KeySyncPerennialStrategy,
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("removing a value from the global cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline:               "0",
				gitconfig.KeySyncPerennialStrategy: "1",
			},
			Local: gitconfig.SingleSnapshot{},
		}
		after := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline: "0",
			},
			Local: gitconfig.SingleSnapshot{},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added: []gitconfig.Key{},
				Removed: map[gitconfig.Key]string{
					gitconfig.KeySyncPerennialStrategy: "1",
				},
				Changed: map[gitconfig.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []gitconfig.Key{},
				Removed: map[gitconfig.Key]string{},
				Changed: map[gitconfig.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.SetGlobalConfig{
				Key:   gitconfig.KeySyncPerennialStrategy,
				Value: "1",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("changing a value in the global cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline: "0",
			},
			Local: gitconfig.SingleSnapshot{},
		}
		after := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline: "1",
			},
			Local: gitconfig.SingleSnapshot{},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []gitconfig.Key{},
				Removed: map[gitconfig.Key]string{},
				Changed: map[gitconfig.Key]undodomain.Change[string]{
					gitconfig.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []gitconfig.Key{},
				Removed: map[gitconfig.Key]string{},
				Changed: map[gitconfig.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.SetGlobalConfig{
				Key:   gitconfig.KeyOffline,
				Value: "0",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("adding a value to the local cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{},
			Local: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline: "0",
			},
		}
		after := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{},
			Local: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline:               "0",
				gitconfig.KeySyncPerennialStrategy: "1",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: emptyConfigDiff(),
			Local: undoconfig.ConfigDiff{
				Added: []gitconfig.Key{
					gitconfig.KeySyncPerennialStrategy,
				},
				Removed: map[gitconfig.Key]string{},
				Changed: map[gitconfig.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.RemoveLocalConfig{
				Key: gitconfig.KeySyncPerennialStrategy,
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("removing a value from the local cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{},
			Local: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline:               "0",
				gitconfig.KeySyncPerennialStrategy: "1",
			},
		}
		after := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{},
			Local: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline: "0",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []gitconfig.Key{},
				Removed: map[gitconfig.Key]string{},
				Changed: map[gitconfig.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added: []gitconfig.Key{},
				Removed: map[gitconfig.Key]string{
					gitconfig.KeySyncPerennialStrategy: "1",
				},
				Changed: map[gitconfig.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcodes.SetLocalConfig{
				Key:   gitconfig.KeySyncPerennialStrategy,
				Value: "1",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("changing a value in the local cache", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{},
			Local: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline: "0",
			},
		}
		after := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{},
			Local: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline: "1",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added:   []gitconfig.Key{},
				Removed: map[gitconfig.Key]string{},
				Changed: map[gitconfig.Key]undodomain.Change[string]{},
			},
			Local: undoconfig.ConfigDiff{
				Added:   []gitconfig.Key{},
				Removed: map[gitconfig.Key]string{},
				Changed: map[gitconfig.Key]undodomain.Change[string]{
					gitconfig.KeyOffline: {
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
				Key:   gitconfig.KeyOffline,
				Value: "0",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("complex example", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline:  "0",
				gitconfig.KeyPushHook: "0",
			},
			Local: gitconfig.SingleSnapshot{
				gitconfig.KeyPerennialBranches: "prod",
				gitconfig.KeyGithubToken:       "token",
			},
		}
		after := undoconfig.ConfigSnapshot{
			Global: gitconfig.SingleSnapshot{
				gitconfig.KeyOffline:               "1",
				gitconfig.KeySyncPerennialStrategy: "1",
			},
			Local: gitconfig.SingleSnapshot{
				gitconfig.KeyPerennialBranches: "prod qa",
				gitconfig.KeyPushHook:          "1",
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undoconfig.ConfigDiff{
				Added: []gitconfig.Key{
					gitconfig.KeySyncPerennialStrategy,
				},
				Removed: map[gitconfig.Key]string{
					gitconfig.KeyPushHook: "0",
				},
				Changed: map[gitconfig.Key]undodomain.Change[string]{
					gitconfig.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: undoconfig.ConfigDiff{
				Added: []gitconfig.Key{
					gitconfig.KeyPushHook,
				},
				Removed: map[gitconfig.Key]string{
					gitconfig.KeyGithubToken: "token",
				},
				Changed: map[gitconfig.Key]undodomain.Change[string]{
					gitconfig.KeyPerennialBranches: {
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
				Key: gitconfig.KeySyncPerennialStrategy,
			},
			&opcodes.SetGlobalConfig{
				Key:   gitconfig.KeyPushHook,
				Value: "0",
			},
			&opcodes.SetGlobalConfig{
				Key:   gitconfig.KeyOffline,
				Value: "0",
			},
			&opcodes.RemoveLocalConfig{
				Key: gitconfig.KeyPushHook,
			},
			&opcodes.SetLocalConfig{
				Key:   gitconfig.KeyGithubToken,
				Value: "token",
			},
			&opcodes.SetLocalConfig{
				Key:   gitconfig.KeyPerennialBranches,
				Value: "prod",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})
}

func emptyConfigDiff() undoconfig.ConfigDiff {
	return undoconfig.ConfigDiff{
		Added:   []gitconfig.Key{},
		Changed: map[gitconfig.Key]undodomain.Change[string]{},
		Removed: map[gitconfig.Key]string{},
	}
}
