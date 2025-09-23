package configdomain

type ProgramFlow int

const (
	ProgramFlowContinue ProgramFlow = 0
	ProgramFlowExit     ProgramFlow = 1
	ProgramFlowRestart  ProgramFlow = 2
)
