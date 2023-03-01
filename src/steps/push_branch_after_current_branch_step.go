package steps

// PushBranchAfterCurrentBranchSteps is a mock step that is used in the undo process
// to push the branch after other steps have been undone.
type PushBranchAfterCurrentBranchSteps struct {
	EmptyStep
}
