package configfile

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/messages"
)

type Config struct {
	Branches
	CodeHosting
	SyncStrategy
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
	FeatureBranches   configdomain.SyncFeatureStrategy
	PerennialBranches configdomain.SyncPerennialStrategy
}

func load() (*Config, error) {
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

func Parse(text string) (*Config, error) {
	var result Config
	_, err := toml.Decode(text, &result)
	if err != nil {
		return nil, fmt.Errorf(messages.ConfigFileWrongInput, ".git-branches.yml", err)
	}
	return &result, err
}
