package configdomain

import (
	"strings"

	"github.com/git-town/git-town/v21/internal/messages"
)

// defines the type of Git configuration used
type ConfigScope string

const (
	ConfigScopeGlobal ConfigScope = "global"
	ConfigScopeLocal  ConfigScope = "local"
)

// GitFlag provides the flag to use when storing configuration data with this scope in Git metadata.
func (self ConfigScope) GitFlag() string {
	switch self {
	case ConfigScopeGlobal:
		return "--global"
	case ConfigScopeLocal:
		return "--local"
	}
	panic(messages.ConfigScopeUnknown)
}

func (self ConfigScope) String() string {
	return string(self)
}

func ParseConfigScope(text string) ConfigScope {
	switch strings.TrimSpace(text) {
	case "local", "":
		return ConfigScopeLocal
	case "global":
		return ConfigScopeGlobal
	default:
		panic(messages.ConfigScopeUnknown)
	}
}
