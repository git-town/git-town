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
	Branches struct {
		Main       string
		Perennials []string
	}
	CodeHosting struct {
		platform       string
		originHostname string
	}
	SyncStrategy struct {
		FeatureBranches   configdomain.SyncFeatureStrategy
		PerennialBranches configdomain.SyncPerennialStrategy
	}
	PushNewbranches        bool
	ShipDeleteRemoteBranch bool
	SyncUpstream           bool
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
	return parse(string(bytes))
}

func parse(text string) (*Config, error) {
	var result Config
	metadata, err := toml.Decode(text, &result)
	fmt.Println("1111111111111111", metadata)
	if err != nil {
		return nil, fmt.Errorf(messages.ConfigFileWrongInput, ".git-branches.yml", err)
	}
	return &result, err
}
