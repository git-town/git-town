package configfile

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v11/src/messages"
)

type ConfigFile struct {
	Branches               Branches
	CodeHosting            CodeHosting  `toml:"code-hosting"`
	SyncStrategy           SyncStrategy `toml:"sync-strategy"`
	PushNewbranches        bool
	ShipDeleteRemoteBranch bool
	SyncUpstream           bool
}

type Branches struct {
	Main       string
	Perennials []string
}

type CodeHosting struct {
	Platform       string
	OriginHostname string
}

type SyncStrategy struct {
	FeatureBranches   string `toml:"feature-branches"`
	PerennialBranches string `toml:"perennial-branches"`
}

func load() (*ConfigFile, error) {
	file, err := os.Open(".git-branches.toml")
	defer file.Close()
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf(messages.ConfigFileCannotRead, ".git-branches.yml", err)
	}
	return Parse(string(bytes))
}

func Parse(text string) (*ConfigFile, error) {
	var result ConfigFile
	_, err := toml.Decode(text, &result)
	if err != nil {
		return nil, fmt.Errorf(messages.ConfigFileWrongInput, ".git-branches.yml", err)
	}
	return &result, err
}
