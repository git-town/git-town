package git

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	Config   RepoConfig
	Backend  BackendCommands
	Frontend FrontendCommands
}

func NewProdRunner(backendRunner BackendRunner, frontendRunner FrontendRunner, config RepoConfig) ProdRunner {
	backendCommands := BackendCommands{
		BackendRunner: backendRunner,
		Config:        &config,
	}
	frontendCommands := FrontendCommands{
		Frontend: frontendRunner,
		Config:   &config,
		Backend:  &backendCommands,
	}
	return ProdRunner{
		Config:   config,
		Backend:  backendCommands,
		Frontend: frontendCommands,
	}
}
