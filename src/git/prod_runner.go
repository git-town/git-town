package git

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	Config      RepoConfig
	Backend     BackendCommands
	Frontend    FrontendCommands
	CommandsRun CommandsRun
}

type CommandsRun interface {
	PrintAnalysis()
}
