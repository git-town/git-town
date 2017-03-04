package steps

type StepList struct {
  List []Step
}

func (stepList *StepList) Append(step Step) {
  stepList.List = append(stepList.List, step)
}

func (stepList *StepList) AppendAll(steps []Step) {
  stepList.List = append(stepList.List, steps...)
}

func (stepList *StepList) AppendList(otherList StepList) {
  stepList.List = append(stepList.List, otherList.List...)
}

func (stepList *StepList) Pop() (result Step) {
  if len(stepList.List) == 0 {
    return nil
  }
  result = stepList.List[0]
  stepList.List = stepList.List[1:]
  return result
}

func (stepList *StepList) Prepend(step Step) {
  stepList.List = append([]Step{step}, stepList.List...)
}

func (stepList *StepList) PrependAll(steps []Step) {
  stepList.List = append(steps, stepList.List...)
}
