package envconfig

import (
	"strings"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// EnvVars is an immutable representation of all environment variables.
// It allows efficient lookup of environment variables in O(1) time
// by multiple names.
type EnvVars struct {
	data map[string]string
}

// Get provides the environment variable with the first matching given name.
// TODO: delete this and use GetOpt everywhere instead. Rename GetOpt to Get.
func (self EnvVars) Get(name string, alternatives ...string) stringss.TrimmedString {
	if result, has := self.data[name]; has {
		return stringss.TrimSpace(result)
	}
	for _, alternative := range alternatives {
		if result, has := self.data[alternative]; has {
			return stringss.TrimSpace(result)
		}
	}
	return ""
}

func (self EnvVars) GetFirstNonEmpty(name string, alternatives ...string) Option[string] {
	if result, has := self.data[name]; has {
		if len(result) > 0 {
			return Some(result)
		}
	}
	for _, alternative := range alternatives {
		if result, has := self.data[alternative]; has {
			if len(result) > 0 {
				return Some(result)
			}
		}
	}
	return None[string]()
}

// GetOpt provides the content of the environment variable with the first matching given name.
func (self EnvVars) GetOpt(name string, alternatives ...string) Option[string] {
	if result, has := self.data[name]; has {
		return Some(result)
	}
	for _, alternative := range alternatives {
		if result, has := self.data[alternative]; has {
			return Some(result)
		}
	}
	return None[string]()
}

func NewEnvVars(entries []string) EnvVars {
	result := EnvVars{
		data: map[string]string{},
	}
	for _, entry := range entries {
		if name, value, isValid := strings.Cut(entry, "="); isValid {
			result.data[name] = value
		}
	}
	return result
}
