package envconfig

import (
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
)

type ImmutableEnvironment map[string]string

func (self ImmutableEnvironment) LoadKey(key configdomain.Key) string {
	return self.LoadString(key.String())
}

func (self ImmutableEnvironment) LoadString(name string) string {
	return self[name]
}

func NewImmutableEnvironment(osEnv []string) ImmutableEnvironment {
	result := ImmutableEnvironment{}
	for _, entry := range osEnv {
		if name, value, isValid := strings.Cut(entry, "="); isValid {
			if key, hasKey := configdomain.ParseKey(name).Get(); hasKey {
				result[key.String()] = value
			}
		}
	}
	return result
}
