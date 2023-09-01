package git

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	Config   RepoConfig
	Backend  BackendCommands
	Frontend FrontendCommands
	Stats    Statistics
}

func EmptyProdRunner() ProdRunner {
	return ProdRunner{
		Config:   RepoConfig{},
		Backend:  BackendCommands{},
		Frontend: FrontendCommands{},
		Stats:    nil,
	}
}

type Statistics interface {
	PrintAnalysis()
}
