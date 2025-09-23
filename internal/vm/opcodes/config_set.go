package opcodes

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

type ConfigSet struct {
	Key   configdomain.Key
	Scope configdomain.ConfigScope
	Value string
}

func (self *ConfigSet) Run(args shared.RunArgs) error {
	return gitconfig.SetConfigValue(args.Backend, self.Scope, self.Key, self.Value)
}
