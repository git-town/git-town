package cli

// Config defines the configuration values needed by the `cli` package.
type Config interface {
	GetBranchAncestryRoots() []string
	GetChildBranches(string) []string
}
