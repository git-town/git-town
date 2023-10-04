package steps_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/shoenig/test/must"
)

func TestStepListBuilder(t *testing.T) {
	t.Parallel()

	t.Run("AppendE", func(t *testing.T) {
		t.Parallel()
		t.Run("without registered error", func(t *testing.T) {
			t.Parallel()
			t.Run("registers the given step", func(t *testing.T) {
				t.Parallel()
				b := steps.StepListBuilder{}
				step := step.Empty{}
				b.AddE(&step, nil)
				must.Eq(t, steps.NewStepList(&step), b.StepList)
			})
			t.Run("registers the given error", func(t *testing.T) {
				t.Parallel()
				b := steps.StepListBuilder{}
				err := errors.New("test error")
				b.AddE(&step.Empty{}, err)
				list, builderErr := b.Result()
				must.True(t, list.IsEmpty())
				must.EqOp(t, err, builderErr)
			})
		})

		t.Run("with an already registered error", func(t *testing.T) {
			t.Parallel()
			t.Run("keeps the already registered error", func(t *testing.T) {
				t.Parallel()
				b := steps.StepListBuilder{}
				firstErr := errors.New("first error")
				b.Check(firstErr)
				b.AddE(&step.Empty{}, errors.New("second error"))
				_, builderErr := b.Result()
				must.EqOp(t, firstErr, builderErr)
			})
			t.Run("does not add the given step", func(t *testing.T) {
				t.Parallel()
				b := steps.StepListBuilder{}
				b.Fail("existing error")
				step := step.Empty{}
				b.AddE(&step, nil)
				list, _ := b.Result()
				must.True(t, list.IsEmpty())
			})
		})
	})
}
