package steps

type StepList struct {
	List []Step
}

func (stepList *StepList) Append(step Step) {
	stepList.List = append(stepList.List, step)
}

func (stepList *StepList) AppendList(otherList StepList) {
	stepList.List = append(stepList.List, otherList.List...)
}

func (stepList *StepList) Peek() (result Step) {
	if len(stepList.List) == 0 {
		return nil
	}
	return stepList.List[0]
}

func (stepList *StepList) Pop() (result Step) {
	if len(stepList.List) == 0 {
		return nil
	}
	result = stepList.List[0]
	stepList.List = stepList.List[1:]
	return
}

func (stepList *StepList) Prepend(step Step) {
	stepList.List = append([]Step{step}, stepList.List...)
}
