package git

// ProdRunner provides Git functionality for production code.
type ProdRunner struct {
	Config   RepoConfig
	Backend  BackendCommands
	Frontend FrontendCommands
}
