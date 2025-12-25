package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v22/internal/messages"
)

// ConfigScope defines the type of Git configuration used.
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

func ParseConfigScope(text string) (ConfigScope, error) {
	switch strings.TrimSpace(text) {
	case "local", "":
		return ConfigScopeLocal, nil
	case "global":
		return ConfigScopeGlobal, nil
	default:
		return "", fmt.Errorf("unknown config scope: %q", text)
	}
}
