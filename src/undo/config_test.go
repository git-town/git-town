package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/shoenig/test/must"
)

func TestConfigUndo(t *testing.T) {
	t.Parallel()

	t.Run("global config added", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline: "0",
				},
				Local: config.GitConfigCache{},
			},
		}
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline:            "0",
					config.KeyPullBranchStrategy: "1",
				},
				Local: config.GitConfigCache{},
			},
		}
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: undo.ConfigDiff{
				Added: []config.Key{
					config.KeyPullBranchStrategy,
				},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
			Local: undo.EmptyConfigDiff(),
		}
		must.Eq(t, wantDiff, haveDiff)
		haveSteps := haveDiff.UndoSteps()
		wantSteps := runstate.StepList{
			List: []step.Step{
				&step.RemoveGlobalConfig{
					Key: config.KeyPullBranchStrategy,
				},
			},
		}
		must.Eq(t, wantSteps, haveSteps)
	})

	t.Run("global config removed", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline:            "0",
					config.KeyPullBranchStrategy: "1",
				},
				Local: config.GitConfigCache{},
			},
		}
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline: "0",
				},
				Local: config.GitConfigCache{},
			},
		}
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: undo.ConfigDiff{
				Added: []config.Key{},
				Removed: map[config.Key]string{
					config.KeyPullBranchStrategy: "1",
				},
				Changed: map[config.Key]domain.Change[string]{},
			},
			Local: undo.ConfigDiff{
				Added:   []config.Key{},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveSteps := haveDiff.UndoSteps()
		wantSteps := runstate.StepList{
			List: []step.Step{
				&step.SetGlobalConfig{
					Key:   config.KeyPullBranchStrategy,
					Value: "1",
				},
			},
		}
		must.Eq(t, wantSteps, haveSteps)
	})

	t.Run("global config changed", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline: "0",
				},
				Local: config.GitConfigCache{},
			},
		}
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{
					config.KeyOffline: "1",
				},
				Local: config.GitConfigCache{},
			},
		}
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: undo.ConfigDiff{
				Added:   []config.Key{},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{
					config.KeyOffline: {
						Before: "0",
						After:  "1",
					},
				},
			},
			Local: undo.ConfigDiff{
				Added:   []config.Key{},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveSteps := haveDiff.UndoSteps()
		wantSteps := runstate.StepList{
			List: []step.Step{
				&step.SetGlobalConfig{
					Key:   config.KeyOffline,
					Value: "0",
				},
			},
		}
		must.Eq(t, wantSteps, haveSteps)
	})

	t.Run("local config added", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline: "0",
				},
			},
		}
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline:            "0",
					config.KeyPullBranchStrategy: "1",
				},
			},
		}
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: undo.EmptyConfigDiff(),
			Local: undo.ConfigDiff{
				Added: []config.Key{
					config.KeyPullBranchStrategy,
				},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveSteps := haveDiff.UndoSteps()
		wantSteps := runstate.StepList{
			List: []step.Step{
				&step.RemoveLocalConfig{
					Key: config.KeyPullBranchStrategy,
				},
			},
		}
		must.Eq(t, wantSteps, haveSteps)
	})

	t.Run("local config removed", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline:            "0",
					config.KeyPullBranchStrategy: "1",
				},
			},
		}
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline: "0",
				},
			},
		}
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: undo.ConfigDiff{
				Added:   []config.Key{},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
			Local: undo.ConfigDiff{
				Added: []config.Key{},
				Removed: map[config.Key]string{
					config.KeyPullBranchStrategy: "1",
				},
				Changed: map[config.Key]domain.Change[string]{},
			},
		}
		must.Eq(t, wantDiff, haveDiff)
		haveSteps := haveDiff.UndoSteps()
		wantSteps := runstate.StepList{
			List: []step.Step{
				&step.SetLocalConfig{
					Key:   config.KeyPullBranchStrategy,
					Value: "1",
				},
			},
		}
		must.Eq(t, wantSteps, haveSteps)
	})

	t.Run("local config changed", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline: "0",
				},
			},
		}
		after := undo.ConfigSnapshot{
			Cwd: "/foo",
			GitConfig: config.GitConfig{
				Global: config.GitConfigCache{},
				Local: config.GitConfigCache{
					config.KeyOffline: "1",
				},
			},
		}
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: undo.ConfigDiff{
				Added:   []config.Key{},
				Removed: map[config.Key]string{},
				Changed: map[config.Key]domain.Change[string]{},
			},
			Local: undo.ConfigDiff{
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
		haveSteps := haveDiff.UndoSteps()
		wantSteps := runstate.StepList{
			List: []step.Step{
				&step.SetLocalConfig{
					Key:   config.KeyOffline,
					Value: "0",
				},
			},
		}
		must.Eq(t, wantSteps, haveSteps)
	})

	t.Run("complex example", func(t *testing.T) {
		t.Parallel()
		before := undo.ConfigSnapshot{
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
		after := undo.ConfigSnapshot{
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
		haveDiff := undo.NewConfigDiffs(before, after)
		wantDiff := undo.ConfigDiffs{
			Global: undo.ConfigDiff{
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
			Local: undo.ConfigDiff{
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
		haveSteps := haveDiff.UndoSteps()
		wantSteps := runstate.StepList{
			List: []step.Step{
				&step.RemoveGlobalConfig{
					Key: config.KeyPullBranchStrategy,
				},
				&step.SetGlobalConfig{
					Key:   config.KeyPushHook,
					Value: "0",
				},
				&step.SetGlobalConfig{
					Key:   config.KeyOffline,
					Value: "0",
				},
				&step.RemoveLocalConfig{
					Key: config.KeyPushHook,
				},
				&step.SetLocalConfig{
					Key:   config.KeyGithubToken,
					Value: "token",
				},
				&step.SetLocalConfig{
					Key:   config.KeyPerennialBranches,
					Value: "prod",
				},
			},
		}
		must.Eq(t, wantSteps, haveSteps)
	})
}
