package cliconfig

import "github.com/git-town/git-town/v21/internal/config/configdomain"

type CliConfig struct {
	DryRun  configdomain.DryRun
	Verbose configdomain.Verbose
}
