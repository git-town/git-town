package steps

// StepList is a list of steps
// with convenience functions for adding and removing steps.
type StepList struct {
	List []Step
}

// Append adds the given step to the end of this StepList.
func (stepList *StepList) Append(step Step) {
	stepList.List = append(stepList.List, step)
}

// AppendList adds all elements of the given StepList to the end of this StepList.
func (stepList *StepList) AppendList(otherList StepList) {
	stepList.List = append(stepList.List, otherList.List...)
}

// Peek returns the first element of this StepList.
func (stepList *StepList) Peek() (result Step) {
	if len(stepList.List) == 0 {
		return nil
	}
	return stepList.List[0]
}

// Pop removes and returns the first element of this StepList.
func (stepList *StepList) Pop() (result Step) {
	if len(stepList.List) == 0 {
		return nil
	}
	result = stepList.List[0]
	stepList.List = stepList.List[1:]
	return
}

// Prepend adds the given step to the beginning of this StepList.
func (stepList *StepList) Prepend(step Step) {
	stepList.List = append([]Step{step}, stepList.List...)
}
