package configfile

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/messages"
)

type ConfigFile struct {
	Branches               Branches
	CodeHosting            CodeHosting  `toml:"code-hosting"`
	SyncStrategy           SyncStrategy `toml:"sync-strategy"`
	PushNewbranches        bool         `toml:"push-new-branches"`
	ShipDeleteRemoteBranch bool         `toml:"ship-delete-remote-branch"`
	SyncUpstream           bool         `toml:"sync-upstream"`
}

type Branches struct {
	Main       string   `toml:"main"`
	Perennials []string `toml:"perennials"`
}

type CodeHosting struct {
	Platform       string
	OriginHostname string
}

type SyncStrategy struct {
	FeatureBranches   string `toml:"feature-branches"`
	PerennialBranches string `toml:"perennial-branches"`
}

func (self SyncStrategy) SyncFeatureStrategy() (configdomain.SyncFeatureStrategy, error) {
	return configdomain.NewSyncFeatureStrategy(self.FeatureBranches)
}

func (self SyncStrategy) SyncPerennialStrategy() (configdomain.SyncPerennialStrategy, error) {
	return configdomain.NewSyncPerennialStrategy(self.PerennialBranches)
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

func Encode(config ConfigFile) string {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(config)
	if err != nil {
		panic(fmt.Sprintf("cannot encode config: %v", err))
	}
	return buf.String()
}

func Parse(text string) (*ConfigFile, error) {
	var result ConfigFile
	_, err := toml.Decode(text, &result)
	if err != nil {
		return nil, fmt.Errorf(messages.ConfigFileWrongInput, ".git-branches.yml", err)
	}
	return &result, err
}
