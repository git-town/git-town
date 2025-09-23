// Package giturl provides facilities to work with the special URL formats used in Git remotes.
package giturl

import (
	"regexp"
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Parts contains recognized parts of a Git URL.
type Parts struct {
	User Option[string] // optional username
	Host string         // hostname of the Git server
	Org  string         // name of the organization that the repo is in
	Repo string         // name of the repository
}

func Parse(url string) Option[Parts] {
	// NOTE: if we can't parse a Git URL, we simply ignore it.
	// This is because the URLs might be on the filesystem.
	// Remotes on the filesystem are not an error condition.

	// handle HTTP/HTTPS URLs
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		httpPattern := regexp.MustCompile(`^https?://(?:([^@]+)@)?([^/]+)/(.*)$`)
		if matches := httpPattern.FindStringSubmatch(url); matches != nil {
			path := strings.TrimSuffix(matches[3], ".git")
			return finalize(matches[1], matches[2], path)
		}
		return None[Parts]()
	}

	// handle SSH URLs with ssh:// prefix
	if strings.HasPrefix(url, "ssh://") {
		sshPattern := regexp.MustCompile(`^ssh://(?:([^@]+)@)?([^/:]+)(?::(\d+))?/(.*)$`)
		if matches := sshPattern.FindStringSubmatch(url); matches != nil {
			path := strings.TrimSuffix(matches[4], ".git")
			return finalize(matches[1], matches[2], path)
		}
		return None[Parts]()
	}

	// handle SSH URLs with colon separator (e.g., git@github.com:user/repo),
	// with and without ports
	colonPattern := regexp.MustCompile(`^(?:([^@]+)@)?([^:]+):(.*)$`)
	if matches := colonPattern.FindStringSubmatch(url); matches != nil {
		host := matches[2]
		path := matches[3]

		// handle port numbers in path (e.g., git@git.example.com:4022/a/b.git)
		if portSlashMatch := regexp.MustCompile(`^(\d+)/(.*)$`).FindStringSubmatch(path); portSlashMatch != nil {
			path = portSlashMatch[2]
		}
		path = strings.TrimSuffix(path, ".git")
		return finalize(matches[1], host, path)
	}

	// handle SSH URLs with slash separator (e.g., git@bitbucket.org/user/repo)
	slashPattern := regexp.MustCompile(`^(?:([^@]+)@)?([^/]+)/(.*)$`)
	if matches := slashPattern.FindStringSubmatch(url); matches != nil {
		path := strings.TrimSuffix(matches[3], ".git")
		return finalize(matches[1], matches[2], path)
	}

	return None[Parts]()
}

func finalize(userMatch, host, path string) Option[Parts] {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return None[Parts]()
	}

	var user Option[string]
	if userMatch != "" {
		user = Some(strings.TrimSuffix(userMatch, "@"))
	}

	var org string
	var repo string

	// Special case for Azure DevOps URLs: remove "v3" prefix from path
	if host == "ssh.dev.azure.com" && len(parts) >= 3 && parts[0] == "v3" {
		parts = parts[1:] // remove the "v3" prefix
	}

	org = strings.Join(parts[:len(parts)-1], "/") // all but the last part are org, last part is repo
	repo = parts[len(parts)-1]

	return Some(Parts{
		Host: host,
		Org:  org,
		Repo: repo,
		User: user,
	})
}
