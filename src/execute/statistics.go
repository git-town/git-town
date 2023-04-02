package execute

type Statistics struct {
	CommandsCount int
}

func (s *Statistics) RegisterCommandExecution() {
	s.CommandsCount += 1
}
