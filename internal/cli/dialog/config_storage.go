package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/messages"
)

const (
	configStorageTitle = `Configuration storage`
	configStorageHelp  = `
How do you want to store the configuration data?

You can store it as a configuration file
(git-town.toml), which you commit
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

func ConfigStorage(inputs dialogcomponents.Inputs) (ConfigStorageOption, dialogdomain.Exit, error) {
	entries := list.NewEntries(
		ConfigStorageOptionGit,
		ConfigStorageOptionFile,
	)
	selection, exit, err := dialogcomponents.RadioList(entries, 0, configStorageTitle, configStorageHelp, inputs, "config-storage")
	fmt.Printf(messages.ConfigStorage, dialogcomponents.FormattedSelection(selection.Short(), exit))
	return selection, exit, err
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
