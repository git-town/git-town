package git

// ProdRepo provides Git functionality for production code.
type ProdRepo struct {
	Config   RepoConfig
	Backend  BackendCommands
	Frontend FrontendCommands
}

func NewProdRepo(omitBranchNames, dryRun, debug bool) ProdRepo {
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
	return ProdRepo{
		Config:   config,
		Backend:  backendCommands,
		Frontend: frontendCommands,
	}
}
