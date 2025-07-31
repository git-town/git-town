package cliconfig

import "github.com/git-town/git-town/v21/internal/config/configdomain"

// CliConfig contains the generic (command-independent)
// configuration information that can be received
// via CLI flags.
type CliConfig struct {
	AutoResolve configdomain.AutoResolve
	DryRun      configdomain.DryRun
	Verbose     configdomain.Verbose
}
