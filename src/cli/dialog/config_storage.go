package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/messages"
)

const (
	configStorageTitle = `Configuration storage`
	configStorageHelp  = `
How do you want to store the configuration data?

You can store it as a configuration file
(.git-branches.toml), which you commit
to the repository. This sets up Git Town
for all people working on this codebase.
Personal data like your API tokens
remain on this machine only.

You can also store the Git Town configuration
as Git metadata on this machine only.

`
)

const (
	ConfigStorageOptionFile ConfigStorageOption = `configuration file`
	ConfigStorageOptionGit  ConfigStorageOption = `Git metadata`
)

func ConfigStorage(hasConfigFile bool, inputs components.TestInput) (ConfigStorageOption, bool, error) {
	if hasConfigFile {
		return ConfigStorageOptionFile, false, nil
	}
	entries := list.NewEntries(
		ConfigStorageOptionFile,
		ConfigStorageOptionGit,
	)
	selection, aborted, err := components.RadioList(entries, 0, configStorageTitle, configStorageHelp, inputs)
	fmt.Printf(messages.ConfigStorage, components.FormattedSelection(selection.Short(), aborted))
	return selection, aborted, err
}

type ConfigStorageOption string

func (self ConfigStorageOption) Short() string {
	switch self {
	case ConfigStorageOptionFile:
		return "file"
	case ConfigStorageOptionGit:
		return "git"
	}
	panic("unhandled config storage option: " + self)
}

func (self ConfigStorageOption) String() string {
	return string(self)
}
