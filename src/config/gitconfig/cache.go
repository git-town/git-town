package gitconfig

import (
	"regexp"
	"sort"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"golang.org/x/exp/maps"
)

type Cache map[configdomain.Key]string

// Clone provides a copy of this GitConfiguration instance.
func (self Cache) Clone() Cache {
	result := Cache{}
	maps.Copy(result, self)
	return result
}

// KeysMatching provides the keys in this GitConfigCache that match the given regex.
func (self Cache) KeysMatching(pattern string) []configdomain.Key {
	result := []configdomain.Key{}
	re := regexp.MustCompile(pattern)
	for key := range self {
		if re.MatchString(key.String()) {
			result = append(result, key)
		}
	}
	sort.Slice(result, func(a, b int) bool { return result[a].String() < result[b].String() })
	return result
}

// LoadGit provides the Git configuration from the given directory or the global one if the global flag is set.
func LoadGitConfigCache(runner Runner, global bool) Cache {
	result := Cache{}
	cmdArgs := []string{"config", "-lz"}
	if global {
		cmdArgs = append(cmdArgs, "--global")
	} else {
		cmdArgs = append(cmdArgs, "--local")
	}
	output, err := runner.Query("git", cmdArgs...)
	if err != nil {
		return result
	}
	if output == "" {
		return result
	}
	for _, line := range strings.Split(output, "\x00") {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "\n", 2)
		key, value := parts[0], parts[1]
		configKey := configdomain.ParseKey(key)
		if configKey != nil {
			result[*configKey] = value
		}
	}
	return result
}
