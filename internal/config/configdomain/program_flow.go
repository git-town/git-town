package configdomain

// ProgramFlow indicates whether to exit or continue the program.
type ProgramFlow int

const (
	ProgramFlowContinue ProgramFlow = 0
	ProgramFlowExit     ProgramFlow = 1
	ProgramFlowRestart  ProgramFlow = 2
)
