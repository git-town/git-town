package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/stretchr/testify/assert"
)

func TestConfigSnapshot(t *testing.T) {
	t.Parallel()

	t.Run("Diff", func(t *testing.T) {
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
			haveDiff := before.Diff(after)
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
			assert.Equal(t, wantDiff, haveDiff)
			haveSteps := haveDiff.UndoSteps()
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.RemoveGlobalConfigStep{
						Key: config.KeyPullBranchStrategy,
					},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
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
			haveDiff := before.Diff(after)
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
			assert.Equal(t, wantDiff, haveDiff)
			haveSteps := haveDiff.UndoSteps()
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.SetGlobalConfigStep{
						Key:   config.KeyPullBranchStrategy,
						Value: "1",
					},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
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
			haveDiff := before.Diff(after)
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
			assert.Equal(t, wantDiff, haveDiff)
			haveSteps := haveDiff.UndoSteps()
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.SetGlobalConfigStep{
						Key:   config.KeyOffline,
						Value: "0",
					},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
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
			haveDiff := before.Diff(after)
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
			assert.Equal(t, wantDiff, haveDiff)
			haveSteps := haveDiff.UndoSteps()
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.RemoveLocalConfigStep{
						Key: config.KeyPullBranchStrategy,
					},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
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
			haveDiff := before.Diff(after)
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
			assert.Equal(t, wantDiff, haveDiff)
			haveSteps := haveDiff.UndoSteps()
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.SetLocalConfigStep{
						Key:   config.KeyPullBranchStrategy,
						Value: "1",
					},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
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
			haveDiff := before.Diff(after)
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
			assert.Equal(t, wantDiff, haveDiff)
			haveSteps := haveDiff.UndoSteps()
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.SetLocalConfigStep{
						Key:   config.KeyOffline,
						Value: "0",
					},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
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
			haveDiff := before.Diff(after)
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
			assert.Equal(t, wantDiff, haveDiff)
			haveSteps := haveDiff.UndoSteps()
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.RemoveGlobalConfigStep{
						Key: config.KeyPullBranchStrategy,
					},
					&steps.SetGlobalConfigStep{
						Key:   config.KeyPushHook,
						Value: "0",
					},
					&steps.SetGlobalConfigStep{
						Key:   config.KeyOffline,
						Value: "0",
					},
					&steps.RemoveLocalConfigStep{
						Key: config.KeyPushHook,
					},
					&steps.SetLocalConfigStep{
						Key:   config.KeyGithubToken,
						Value: "token",
					},
					&steps.SetLocalConfigStep{
						Key:   config.KeyPerennialBranches,
						Value: "prod",
					},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})
	})
}
