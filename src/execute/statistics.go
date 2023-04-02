package execute

type Statistics struct {
	CommandsCount int
}

func (s *Statistics) RegisterRun() {
	s.CommandsCount += 1
}

func (s *Statistics) RunCount() int {
	return s.CommandsCount
}
