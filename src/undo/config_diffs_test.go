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

func TestConfigDiffs(t *testing.T) {
	t.Parallel()
	t.Run("UndoSteps", func(t *testing.T) {
		t.Parallel()
		t.Run("global config removed", func(t *testing.T) {
			t.Parallel()
			diff := undo.ConfigDiffs{
				Global: undo.ConfigDiff{
					Added: []config.Key{},
					Removed: map[config.Key]string{
						config.KeyOffline: "1",
					},
					Changed: map[config.Key]domain.Change[string]{},
				},
				Local: undo.EmptyConfigDiff(),
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
			t.Parallel()
			diff := undo.ConfigDiffs{
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
				Local: undo.EmptyConfigDiff(),
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
			t.Parallel()
			diff := undo.ConfigDiffs{
				Global: undo.EmptyConfigDiff(),
				Local: undo.ConfigDiff{
					Added: []config.Key{
						config.KeyOffline,
					},
					Removed: map[config.Key]string{},
					Changed: map[config.Key]domain.Change[string]{},
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
			t.Parallel()
			diff := undo.ConfigDiffs{
				Global: undo.EmptyConfigDiff(),
				Local: undo.ConfigDiff{
					Added: []config.Key{},
					Removed: map[config.Key]string{
						config.KeyOffline: "1",
					},
					Changed: map[config.Key]domain.Change[string]{},
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

		t.Run("local config changed", func(t *testing.T) {
			t.Parallel()
			diff := undo.ConfigDiffs{
				Global: undo.EmptyConfigDiff(),
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
			have := diff.UndoSteps()
			want := runstate.StepList{
				List: []steps.Step{
					&steps.SetLocalConfigStep{
						Key:   config.KeyOffline,
						Value: "0",
					},
				},
			}
			assert.Equal(t, want, have)
		})
	})
}
