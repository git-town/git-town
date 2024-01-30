package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/cli/dialog/components"
)

const (
	configStorageTitle = `Configuration storage`
	configStorageHelp  = `
You have two options to store your configuration data.

You can store it as a configuration file (.git-branches.toml)
and commit it to your repository. This shares the
Git Town configuration with your team, i.e. your team can
use Git Town without having to run the setup assistant.
Personal data like your API tokens remain on this machine only.

If you cannot or don't want to commit another configuration file,
you can store the Git Town configuration as Git metadata.
This way all of it remains on this machine and you don't have
to commit anything.

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
	entries := []ConfigStorageOption{
		ConfigStorageOptionFile,
		ConfigStorageOptionGit,
	}
	selection, aborted, err := components.RadioList(entries, 0, configStorageTitle, configStorageHelp, inputs)
	fmt.Printf("Config storage: %s\n", components.FormattedSelection(selection.Short(), aborted))
	return selection, aborted, err
}

type ConfigStorageOption string

func (self ConfigStorageOption) Short() string {
	start, _, _ := strings.Cut(self.String(), " ")
	return start
}

func (self ConfigStorageOption) String() string {
	return string(self)
}
