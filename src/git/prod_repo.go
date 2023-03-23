package git

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	Config   RepoConfig
	Backend  BackendCommands
	Frontend FrontendCommands
}

func NewProdRunner(omitBranchNames, dryRun, debug bool) ProdRunner {
	backendRunner := NewBackendRunner(nil, debug)
	config := NewRepoConfig(backendRunner)
	backendCommands := BackendCommands{
		BackendRunner: backendRunner,
		Config:        &config,
	}
	frontendRunner := NewFrontendRunner(omitBranchNames, dryRun, config.CurrentBranchCache)
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
