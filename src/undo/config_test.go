package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/undo"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/shoenig/test/must"
)

func TestConfigUndo(t *testing.T) {
	t.Parallel()

	t.Run("global cache added", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				GlobalConfig: configdomain.PartialConfig{}, //nolint:exhaustruct
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.PartialConfig{}, //nolint:exhaustruct
			},
		}
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
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
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: configdomain.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeySyncPerennialStrategy,
				},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			},
			Local: configdomain.EmptyConfigDiff(),
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

	t.Run("global config added", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.PartialConfig{}, //nolint:exhaustruct
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.PartialConfig{}, //nolint:exhaustruct
			},
		}
		perennialsAfter := domain.NewLocalBranchNames("one", "two")
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{},
				GlobalConfig: configdomain.PartialConfig{ //nolint:exhaustruct
					PerennialBranches: &perennialsAfter,
				},
				LocalCache:  gitconfig.SingleCache{},
				LocalConfig: configdomain.PartialConfig{}, //nolint:exhaustruct
			},
		}
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: configdomain.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeyPerennialBranches,
				},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			},
			Local: configdomain.EmptyConfigDiff(),
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

	t.Run("global config removed", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
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
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				GlobalConfig: configdomain.PartialConfig{}, //nolint:exhaustruct
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.PartialConfig{}, //nolint:exhaustruct
			},
		}
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: configdomain.ConfigDiff{
				Added: []configdomain.Key{},
				Removed: map[configdomain.Key]string{
					configdomain.KeySyncPerennialStrategy: "1",
				},
				Changed: map[configdomain.Key]domain.Change[string]{},
			},
			Local: configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
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

	t.Run("global config changed", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.EmptyPartialConfig(),
			},
		}
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: gitconfig.FullCache{
				GlobalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "1",
				},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache:   gitconfig.SingleCache{},
				LocalConfig:  configdomain.EmptyPartialConfig(),
			},
		}
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
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

	t.Run("local config added", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
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
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: configdomain.EmptyConfigDiff(),
			Local: configdomain.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeySyncPerennialStrategy,
				},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
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

	t.Run("local config removed", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
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
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			},
			Local: configdomain.ConfigDiff{
				Added: []configdomain.Key{},
				Removed: map[configdomain.Key]string{
					configdomain.KeySyncPerennialStrategy: "1",
				},
				Changed: map[configdomain.Key]domain.Change[string]{},
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

	t.Run("local config changed", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "0",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: gitconfig.FullCache{
				GlobalCache:  gitconfig.SingleCache{},
				GlobalConfig: configdomain.EmptyPartialConfig(),
				LocalCache: gitconfig.SingleCache{
					configdomain.KeyOffline: "1",
				},
				LocalConfig: configdomain.EmptyPartialConfig(),
			},
		}
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			},
			Local: configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{
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
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
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
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
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
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: configdomain.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeySyncPerennialStrategy,
				},
				Removed: map[configdomain.Key]string{
					configdomain.KeyPushHook: "0",
				},
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: configdomain.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeyPushHook,
				},
				Removed: map[configdomain.Key]string{
					configdomain.KeyGithubToken: "token",
				},
				Changed: map[configdomain.Key]domain.Change[string]{
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
