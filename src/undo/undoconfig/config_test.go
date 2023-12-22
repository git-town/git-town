package undoconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/undo/undoconfig"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/shoenig/test/must"
)

func TestConfigUndo(t *testing.T) {
	t.Parallel()

	t.Run("adding a value to the global cache", func(t *testing.T) {
		t.Parallel()
		before := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				GlobalConfig: configdomain.PartialConfig{}, //nolint:exhaustruct
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.PartialConfig{}, //nolint:exhaustruct
			},
		}
		after := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline:               "0",
					configdomain.KeySyncPerennialStrategy: "1",
				},
				GlobalConfig: configdomain.PartialConfig{}, //nolint:exhaustruct
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.PartialConfig{}, //nolint:exhaustruct
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undodomain.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeySyncPerennialStrategy,
				},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
			Local: undodomain.EmptyConfigDiff(),
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcode.RemoveGlobalConfig{
				Key: configdomain.KeySyncPerennialStrategy,
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("adding a value to the global config", func(t *testing.T) {
		t.Parallel()
		before := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.PartialConfig{}, //nolint:exhaustruct
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.PartialConfig{}, //nolint:exhaustruct
			},
		}
		perennialsAfter := gitdomain.NewLocalBranchNames("one", "two")
		after := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{},
				GlobalConfig: configdomain.PartialConfig{ //nolint:exhaustruct
					PerennialBranches: &perennialsAfter,
				},
				LocalCache:  gitconfig.SingleCache{},
				LocalConfig: configdomain.PartialConfig{}, //nolint:exhaustruct
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undodomain.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeyPerennialBranches,
				},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
			Local: undodomain.EmptyConfigDiff(),
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcode.RemoveGlobalConfig{
				Key: configdomain.KeyPerennialBranches,
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("removing a value from the global cache", func(t *testing.T) {
		t.Parallel()
		before := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline:               "0",
					configdomain.KeySyncPerennialStrategy: "1",
				},
				GlobalConfig: configdomain.PartialConfig{}, //nolint:exhaustruct
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.PartialConfig{}, //nolint:exhaustruct
			},
		}
		after := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				GlobalConfig: configdomain.PartialConfig{}, //nolint:exhaustruct
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.PartialConfig{}, //nolint:exhaustruct
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undodomain.ConfigDiff{
				Added: []configdomain.Key{},
				Removed: map[configdomain.Key]string{
					configdomain.KeySyncPerennialStrategy: "1",
				},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
			Local: undodomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcode.SetGlobalConfig{
				Key:   configdomain.KeySyncPerennialStrategy,
				Value: "1",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("changing a value in the global cache", func(t *testing.T) {
		t.Parallel()
		before := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.EmptyPartialConfig(),
			},
		}
		after := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "1",
				},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.EmptyPartialConfig(),
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undodomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: undodomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcode.SetGlobalConfig{
				Key:   configdomain.KeyOffline,
				Value: "0",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("adding a value to the local cache", func(t *testing.T) {
		t.Parallel()
		before := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		after := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyOffline:               "0",
					configdomain.KeySyncPerennialStrategy: "1",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undodomain.EmptyConfigDiff(),
			Local: undodomain.ConfigDiff{
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
			&opcode.RemoveLocalConfig{
				Key: configdomain.KeySyncPerennialStrategy,
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("removing a value from the local cache", func(t *testing.T) {
		t.Parallel()
		before := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyOffline:               "0",
					configdomain.KeySyncPerennialStrategy: "1",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		after := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undodomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
			Local: undodomain.ConfigDiff{
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
			&opcode.SetLocalConfig{
				Key:   configdomain.KeySyncPerennialStrategy,
				Value: "1",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("changing a value in the local cache", func(t *testing.T) {
		t.Parallel()
		before := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		after := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "1",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undodomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			},
			Local: undodomain.ConfigDiff{
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
			&opcode.SetLocalConfig{
				Key:   configdomain.KeyOffline,
				Value: "0",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("complex example", func(t *testing.T) {
		t.Parallel()
		before := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline:  "0",
					configdomain.KeyPushHook: "0",
				},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyPerennialBranches: "prod",
					configdomain.KeyGithubToken:       "token",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		after := undodomain.ConfigSnapshot{
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline:               "1",
					configdomain.KeySyncPerennialStrategy: "1",
				},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyPerennialBranches: "prod qa",
					configdomain.KeyPushHook:          "1",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		haveDiff := undoconfig.NewConfigDiffs(before, after)
		wantDiff := undoconfig.ConfigDiffs{
			Global: undodomain.ConfigDiff{
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
			Local: undodomain.ConfigDiff{
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
			&opcode.RemoveGlobalConfig{
				Key: configdomain.KeySyncPerennialStrategy,
			},
			&opcode.SetGlobalConfig{
				Key:   configdomain.KeyPushHook,
				Value: "0",
			},
			&opcode.SetGlobalConfig{
				Key:   configdomain.KeyOffline,
				Value: "0",
			},
			&opcode.RemoveLocalConfig{
				Key: configdomain.KeyPushHook,
			},
			&opcode.SetLocalConfig{
				Key:   configdomain.KeyGithubToken,
				Value: "token",
			},
			&opcode.SetLocalConfig{
				Key:   configdomain.KeyPerennialBranches,
				Value: "prod",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})
}
