package undoconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/undo/undoconfig"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/program"
	"github.com/shoenig/test/must"
)

func TestConfigUndo(t *testing.T) {
	t.Parallel()

	t.Run("global config added", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline: "0",
				},
				Local: config.GitConfigCache{},
			},
		}
		after := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline:            "0",
					config.KeyPullBranchStrategy: "1",
				},
				Local: config.GitConfigCache{},
			},
		}
		haveDiff := undoconfig.NewDiffs(before, after)
		wantDiff := undoconfig.Diffs{
			Global: undoconfig.Diff{
				Added: []config.Key{
					config.KeyPullBranchStrategy,
				},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
			Local: undoconfig.EmptyDiff(),
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcode.RemoveGlobalConfig{
				Key: config.KeyPullBranchStrategy,
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("global config removed", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline:            "0",
					config.KeyPullBranchStrategy: "1",
				},
				Local: config.GitConfigCache{},
			},
		}
		after := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline: "0",
				},
				Local: config.GitConfigCache{},
			},
		}
		haveDiff := undoconfig.NewDiffs(before, after)
		wantDiff := undoconfig.Diffs{
			Global: undoconfig.Diff{
				Added: []config.Key{},
				Removed: map[config.Key]string{
					config.KeyPullBranchStrategy: "1",
				},
				Changed: map[config.Key]domain.Change[string]{},
			},
			Local: undoconfig.Diff{
				Added:   []config.Key{},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcode.SetGlobalConfig{
				Key:   config.KeyPullBranchStrategy,
				Value: "1",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("global config changed", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline: "0",
				},
				Local: config.GitConfigCache{},
			},
		}
		after := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline: "1",
				},
				Local: config.GitConfigCache{},
			},
		}
		haveDiff := undoconfig.NewDiffs(before, after)
		wantDiff := undoconfig.Diffs{
			Global: undoconfig.Diff{
				Added:   []config.Key{},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{
					config.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: undoconfig.Diff{
				Added:   []config.Key{},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcode.SetGlobalConfig{
				Key:   config.KeyOffline,
				Value: "0",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local config added", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline: "0",
				},
			},
		}
		after := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline:            "0",
					config.KeyPullBranchStrategy: "1",
				},
			},
		}
		haveDiff := undoconfig.NewDiffs(before, after)
		wantDiff := undoconfig.Diffs{
			Global: undoconfig.EmptyDiff(),
			Local: undoconfig.Diff{
				Added: []config.Key{
					config.KeyPullBranchStrategy,
				},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcode.RemoveLocalConfig{
				Key: config.KeyPullBranchStrategy,
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local config removed", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline:            "0",
					config.KeyPullBranchStrategy: "1",
				},
			},
		}
		after := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline: "0",
				},
			},
		}
		haveDiff := undoconfig.NewDiffs(before, after)
		wantDiff := undoconfig.Diffs{
			Global: undoconfig.Diff{
				Added:   []config.Key{},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
			Local: undoconfig.Diff{
				Added: []config.Key{},
				Removed: map[config.Key]string{
					config.KeyPullBranchStrategy: "1",
				},
				Changed: map[config.Key]domain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveProgram := haveDiff.UndoProgram()
		wantProgram := program.Program{
			&opcode.SetLocalConfig{
				Key:   config.KeyPullBranchStrategy,
				Value: "1",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local config changed", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline: "0",
				},
			},
		}
		after := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline: "1",
				},
			},
		}
		haveDiff := undoconfig.NewDiffs(before, after)
		wantDiff := undoconfig.Diffs{
			Global: undoconfig.Diff{
				Added:   []config.Key{},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
			Local: undoconfig.Diff{
				Added:   []config.Key{},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{
					config.KeyOffline: {
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
				Key:   config.KeyOffline,
				Value: "0",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("complex example", func(t *testing.T) {
		t.Parallel()
		before := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline:  "0",
					config.KeyPushHook: "0",
				},
				Local: config.GitConfigCache{
					config.KeyPerennialBranches: "prod",
					config.KeyGithubToken:       "token",
				},
			},
		}
		after := undoconfig.Snapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline:            "1",
					config.KeyPullBranchStrategy: "1",
				},
				Local: config.GitConfigCache{
					config.KeyPerennialBranches: "prod qa",
					config.KeyPushHook:          "1",
				},
			},
		}
		haveDiff := undoconfig.NewDiffs(before, after)
		wantDiff := undoconfig.Diffs{
			Global: undoconfig.Diff{
				Added: []config.Key{
					config.KeyPullBranchStrategy,
				},
				Removed: map[config.Key]string{
					config.KeyPushHook: "0",
				},
				Changed: map[config.Key]domain.Change[string]{
					config.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: undoconfig.Diff{
				Added: []config.Key{
					config.KeyPushHook,
				},
				Removed: map[config.Key]string{
					config.KeyGithubToken: "token",
				},
				Changed: map[config.Key]domain.Change[string]{
					config.KeyPerennialBranches: {
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
				Key: config.KeyPullBranchStrategy,
			},
			&opcode.SetGlobalConfig{
				Key:   config.KeyPushHook,
				Value: "0",
			},
			&opcode.SetGlobalConfig{
				Key:   config.KeyOffline,
				Value: "0",
			},
			&opcode.RemoveLocalConfig{
				Key: config.KeyPushHook,
			},
			&opcode.SetLocalConfig{
				Key:   config.KeyGithubToken,
				Value: "token",
			},
			&opcode.SetLocalConfig{
				Key:   config.KeyPerennialBranches,
				Value: "prod",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})
}
