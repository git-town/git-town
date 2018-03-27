package steps_test

import (
	"encoding/json"

	"github.com/Originate/git-town/src/steps"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RunState", func() {
	Describe("MarshalJSON / UnmarshalJSON", func() {
		It("keeps the step data", func() {
			runState := &steps.RunState{
				AbortStepList: steps.StepList{
					List: []steps.Step{&steps.ResetToShaStep{Sha: "abc"}},
				},
				Command: "sync",
				RunStepList: steps.StepList{
					List: []steps.Step{&steps.ResetToShaStep{Sha: "abc"}},
				},
				UndoStepList: steps.StepList{
					List: []steps.Step{&steps.ResetToShaStep{Sha: "abc"}},
				},
			}
			data, err := json.Marshal(runState)
			Expect(err).To(BeNil())
			newRunState := &steps.RunState{}
			err = json.Unmarshal(data, &newRunState)
			Expect(err).To(BeNil())
			Expect(newRunState).To(Equal(runState))
		})
	})
})
