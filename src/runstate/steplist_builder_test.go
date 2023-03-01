package runstate_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/stretchr/testify/assert"
)

func TestStepListBuilder(t *testing.T) {
	t.Parallel()
	t.Run("AppendE", func(t *testing.T) {
		t.Parallel()
		t.Run("without registered error", func(t *testing.T) {
			t.Parallel()
			t.Run("registers the given step", func(t *testing.T) {
				t.Parallel()
				b := runstate.StepListBuilder{}
				step := steps.EmptyStep{}
				b.AddE(&step, nil)
				assert.Equal(t, runstate.NewStepList(&step), b.StepList)
			})
			t.Run("registers the given error", func(t *testing.T) {
				t.Parallel()
				b := runstate.StepListBuilder{}
				err := errors.New("test error")
				b.AddE(&steps.EmptyStep{}, err)
				list, builderErr := b.Result()
				assert.True(t, list.IsEmpty())
				assert.Equal(t, err, builderErr)
			})
		})

		t.Run("with an already registered error", func(t *testing.T) {
			t.Parallel()
			t.Run("keeps the already registered error", func(t *testing.T) {
				t.Parallel()
				b := runstate.StepListBuilder{}
				firstErr := errors.New("first error")
				b.Check(firstErr)
				b.AddE(&steps.EmptyStep{}, errors.New("second error"))
				_, builderErr := b.Result()
				assert.Error(t, firstErr, builderErr)
			})
			t.Run("does not add the given step", func(t *testing.T) {
				t.Parallel()
				b := runstate.StepListBuilder{}
				b.Fail("existing error")
				step := steps.EmptyStep{}
				b.AddE(&step, nil)
				list, _ := b.Result()
				assert.True(t, list.IsEmpty())
			})
		})
	})
}
