package opcodes

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

type ConfigRemove struct {
	Key                     configdomain.Key // the config key to remove
	Scope                   configdomain.ConfigScope
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConfigRemove) Run(args shared.RunArgs) error {
	return gitconfig.RemoveConfigValue(args.Backend, self.Scope, self.Key)
}
