package configdomain

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

// UserConfigDir is the directory that contains the user-specific configuration on this machine,
// i.e. ~/.config.
type UserConfigDir string

// RepoConfigDir provides the file path where Git Town stores data for the given Git repo inside this UserConfigDir.
func (self UserConfigDir) RepoConfigDir(repoDir gitdomain.RepoRootDir) RepoConfigDir {
	return RepoConfigDir(filepath.Join(self.String(), "git-town", SanitizePath(repoDir.String())))
}

func (self UserConfigDir) String() string {
	return string(self)
}

// TODO: as part of separating low-level (system) and high-level (Git Town specific) code,
// move this and SystemUserConfigDir into package internal/sys/files?
// low-level code cannot depend on high-level code
//
// possible layout:
//
// # DOMAIN LEVEL: Git Town specific things that use the mid and low level
//
// config  execute  proposallineage  setup  skip  state  undo  vm
//
// # MID LEVEL: generic frameworks created or used, uses the low level
//
// cmd  messages  validate
//
// # LOW LEVEL: helpers for interacting with the system
//
// browser  cli  forge  git  gohacks  regexes  subshell
//
// a
func SanitizePath[T ~string](dir T) T { //nolint:ireturn
	replaceCharacterRE := regexp.MustCompile("[[:^alnum:]]")
	sanitized := replaceCharacterRE.ReplaceAllString(string(dir), "-")
	sanitized = strings.ToLower(sanitized)
	replaceDoubleMinusRE := regexp.MustCompile("--+") // two or more dashes
	sanitized = replaceDoubleMinusRE.ReplaceAllString(sanitized, "-")
	for strings.HasPrefix(sanitized, "-") {
		sanitized = sanitized[1:]
	}
	return T(sanitized)
}
