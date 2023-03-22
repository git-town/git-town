package git

// ProdRepo provides Git functionality for production code.
type ProdRepo struct {
	Config   RepoConfig
	Backend  BackendCommands
	Frontend FrontendCommands
}
