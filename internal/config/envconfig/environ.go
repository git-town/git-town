package envconfig

import (
	"fmt"
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
		if envName, value, isValid := strings.Cut(entry, "="); isValid {
			keyName := Env2Key(envName)
			fmt.Println("1111111111111111111111111111111", keyName)
			result[keyName] = value
		}
	}
	return result
}

func Env2Key(envName string) string {
	result := strings.ToLower(envName)
	result = strings.Replace(result, "GIT_TOWN_", "git-town.", 1)
	result = strings.ReplaceAll(result, "_", "-")
	return result
}
