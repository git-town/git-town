package git

// ProdRepo provides Git functionality for production code.
type ProdRepo struct {
	Config   RepoConfig
	Internal InternalCommands
	Public   PublicCommands
}
