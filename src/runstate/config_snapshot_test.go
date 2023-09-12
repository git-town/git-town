package runstate_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/stretchr/testify/assert"
)

func TestConfigSnapshot(t *testing.T) {
	t.Parallel()
	t.Run("Diff", func(t *testing.T) {

		t.Run("global config added", func(t *testing.T) {
			before := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{
						config.KeyOffline: "0",
					},
					Local: config.GitConfigCache{},
				},
			}
			after := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{
						config.KeyOffline:            "0",
						config.KeyPullBranchStrategy: "1",
					},
					Local: config.GitConfigCache{},
				},
			}
			have := before.Diff(after)
			want := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{
					Added: []config.Key{
						config.KeyPullBranchStrategy,
					},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{},
				},
				Local: runstate.ConfigDiff{
					Added:   []config.Key{},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{},
				},
			}
			assert.Equal(t, want, have)
		})

		t.Run("global config removed", func(t *testing.T) {
			before := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{
						config.KeyOffline:            "0",
						config.KeyPullBranchStrategy: "1",
					},
					Local: config.GitConfigCache{},
				},
			}
			after := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{
						config.KeyOffline: "0",
					},
					Local: config.GitConfigCache{},
				},
			}
			have := before.Diff(after)
			want := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{
					Added: []config.Key{},
					Removed: map[config.Key]string{
						config.KeyPullBranchStrategy: "1",
					},
					Changed: map[config.Key]runstate.Change[string]{},
				},
				Local: runstate.ConfigDiff{
					Added:   []config.Key{},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{},
				},
			}
			assert.Equal(t, want, have)
		})

		t.Run("global config changed", func(t *testing.T) {
			before := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{
						config.KeyOffline: "0",
					},
					Local: config.GitConfigCache{},
				},
			}
			after := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{
						config.KeyOffline: "1",
					},
					Local: config.GitConfigCache{},
				},
			}
			have := before.Diff(after)
			want := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{
					Added:   []config.Key{},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{
						config.KeyOffline: {
							Before: "0",
							After:  "1",
						},
					},
				},
				Local: runstate.ConfigDiff{
					Added:   []config.Key{},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{},
				},
			}
			assert.Equal(t, want, have)
		})

		t.Run("local config added", func(t *testing.T) {
			before := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{},
					Local: config.GitConfigCache{
						config.KeyOffline: "0",
					},
				},
			}
			after := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{},
					Local: config.GitConfigCache{
						config.KeyOffline:            "0",
						config.KeyPullBranchStrategy: "1",
					},
				},
			}
			have := before.Diff(after)
			want := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{
					Added:   []config.Key{},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{},
				},
				Local: runstate.ConfigDiff{
					Added: []config.Key{
						config.KeyPullBranchStrategy,
					},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{},
				},
			}
			assert.Equal(t, want, have)
		})

		t.Run("local config removed", func(t *testing.T) {
			before := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{},
					Local: config.GitConfigCache{
						config.KeyOffline:            "0",
						config.KeyPullBranchStrategy: "1",
					},
				},
			}
			after := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{},
					Local: config.GitConfigCache{
						config.KeyOffline: "0",
					},
				},
			}
			have := before.Diff(after)
			want := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{
					Added:   []config.Key{},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{},
				},
				Local: runstate.ConfigDiff{
					Added: []config.Key{},
					Removed: map[config.Key]string{
						config.KeyPullBranchStrategy: "1",
					},
					Changed: map[config.Key]runstate.Change[string]{},
				},
			}
			assert.Equal(t, want, have)
		})

		t.Run("local config changed", func(t *testing.T) {
			before := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{},
					Local: config.GitConfigCache{
						config.KeyOffline: "0",
					},
				},
			}
			after := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{},
					Local: config.GitConfigCache{
						config.KeyOffline: "1",
					},
				},
			}
			have := before.Diff(after)
			want := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{
					Added:   []config.Key{},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{},
				},
				Local: runstate.ConfigDiff{
					Added:   []config.Key{},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{
						config.KeyOffline: {
							Before: "0",
							After:  "1",
						},
					},
				},
			}
			assert.Equal(t, want, have)
		})

		t.Run("complex example", func(t *testing.T) {
			before := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{
						config.KeyOffline:  "0",
						config.KeyPushHook: "0",
					},
					Local: config.GitConfigCache{
						config.KeyMainBranch:        "main",
						config.KeyPerennialBranches: "prod",
					},
				},
			}
			after := runstate.ConfigSnapshot{
				Cwd: "/foo",
				GitConfig: config.GitConfig{
					Global: config.GitConfigCache{
						config.KeyOffline:            "1",
						config.KeyPullBranchStrategy: "1",
					},
					Local: config.GitConfigCache{
						config.KeyMainBranch:        "dev",
						config.KeyPerennialBranches: "prod qa",
						config.KeyPushHook:          "1",
					},
				},
			}
			have := before.Diff(after)
			want := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{
					Added: []config.Key{
						config.KeyPullBranchStrategy,
					},
					Removed: map[config.Key]string{
						config.KeyPushHook: "0",
					},
					Changed: map[config.Key]runstate.Change[string]{
						config.KeyOffline: {
							Before: "0",
							After:  "1",
						},
					},
				},
				Local: runstate.ConfigDiff{
					Added: []config.Key{
						config.KeyPushHook,
					},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{
						config.KeyMainBranch: {
							Before: "main",
							After:  "dev",
						},
						config.KeyPerennialBranches: {
							Before: "prod",
							After:  "prod qa",
						},
					},
				},
			}
			assert.Equal(t, want, have)
		})
	})
}

func TestSnapshotConfigDiff(t *testing.T) {
	t.Parallel()
	t.Run("UndoSteps", func(t *testing.T) {
		t.Parallel()
		t.Run("global config added", func(t *testing.T) {
			diff := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{
					Added: []config.Key{
						config.KeyOffline,
					},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{},
				},
				Local: runstate.ConfigDiff{},
			}
			have := diff.UndoSteps()
			want := runstate.StepList{
				List: []steps.Step{
					&steps.RemoveGlobalConfigStep{
						Key: config.KeyOffline,
					},
				},
			}
			assert.Equal(t, want, have)
		})

		t.Run("global config removed", func(t *testing.T) {
			diff := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{
					Added: []config.Key{},
					Removed: map[config.Key]string{
						config.KeyOffline: "1",
					},
					Changed: map[config.Key]runstate.Change[string]{},
				},
				Local: runstate.ConfigDiff{},
			}
			have := diff.UndoSteps()
			want := runstate.StepList{
				List: []steps.Step{
					&steps.SetGlobalConfigStep{
						Key:   config.KeyOffline,
						Value: "1",
					},
				},
			}
			assert.Equal(t, want, have)
		})

		t.Run("global config changed", func(t *testing.T) {
			diff := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{
					Added:   []config.Key{},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{
						config.KeyOffline: {
							Before: "0",
							After:  "1",
						},
					},
				},
				Local: runstate.ConfigDiff{},
			}
			have := diff.UndoSteps()
			want := runstate.StepList{
				List: []steps.Step{
					&steps.SetGlobalConfigStep{
						Key:   config.KeyOffline,
						Value: "0",
					},
				},
			}
			assert.Equal(t, want, have)
		})

		t.Run("local config added", func(t *testing.T) {
			diff := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{},
				Local: runstate.ConfigDiff{
					Added: []config.Key{
						config.KeyOffline,
					},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]runstate.Change[string]{},
				},
			}
			have := diff.UndoSteps()
			want := runstate.StepList{
				List: []steps.Step{
					&steps.RemoveLocalConfigStep{
						Key: config.KeyOffline,
					},
				},
			}
			assert.Equal(t, want, have)
		})

		t.Run("local config removed", func(t *testing.T) {
			diff := runstate.SnapshotConfigDiff{
				Global: runstate.ConfigDiff{},
				Local: runstate.ConfigDiff{
					Added: []config.Key{},
					Removed: map[config.Key]string{
						config.KeyOffline: "1",
					},
					Changed: map[config.Key]runstate.Change[string]{},
				},
			}
			have := diff.UndoSteps()
			want := runstate.StepList{
				List: []steps.Step{
					&steps.SetLocalConfigStep{
						Key:   config.KeyOffline,
						Value: "1",
					},
				},
			}
			assert.Equal(t, want, have)
		})
	})
}
