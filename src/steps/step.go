package steps

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Originate/git-town/src/git"
)

// Step represents a dedicated activity within a Git Town command.
// Git Town commands are comprised of a number of steps that need to be executed.
type Step interface {
	CreateAbortStep() Step
	CreateContinueStep() Step
	CreateUndoStepBeforeRun() Step
	CreateUndoStepAfterRun() Step
	GetAutomaticAbortErrorMessage() string
	Run() error
	ShouldAutomaticallyAbortOnError() bool
}

// SerializedStep is used to store Steps as JSON.
type SerializedStep struct {
	Data []byte
	Type string
}

// SerializedRunState is used to store RunStates as JSON.
type SerializedRunState struct {
	AbortStep SerializedStep
	RunSteps  []SerializedStep
	UndoSteps []SerializedStep
}

func getRunResultFilename(command string) string {
	replaceCharacterRegexp, err := regexp.Compile("[[:^alnum:]]")
	if err != nil {
		log.Fatal("Error compiling replace character expression: ", err)
	}
	directory := replaceCharacterRegexp.ReplaceAllString(git.GetRootDirectory(), "-")
	return fmt.Sprintf("/tmp/%s_%s", command, directory)
}
