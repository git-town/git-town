package config

import (
	"regexp"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

type GitConfigCache map[Key]string

// Clone provides a copy of this GitConfiguration instance.
func (gcc GitConfigCache) Clone() GitConfigCache {
	result := GitConfigCache{}
	maps.Copy(result, gcc)
	return result
}

// KeysMatching provides the keys in this GitConfigCache that match the given regex.
func (gcc GitConfigCache) KeysMatching(pattern string) []Key {
	result := []Key{}
	re := regexp.MustCompile(pattern)
	for key := range gcc {
		if re.MatchString(key.String()) {
			result = append(result, key)
		}
	}
	sort.Slice(result, func(a, b int) bool { return result[a].Name < result[b].Name })
	return result
}

// LoadGit provides the Git configuration from the given directory or the global one if the global flag is set.
func LoadGitConfigCache(runner runner, global bool) GitConfigCache {
	result := GitConfigCache{}
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
		configKey := ParseKey(key)
		if configKey != nil {
			result[*configKey] = value
		}
	}
	return result
}
