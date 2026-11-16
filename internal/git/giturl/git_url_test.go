package giturl_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/git/giturl"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParse(t *testing.T) {
	t.Parallel()
	tests := map[string]giturl.Parts{
		"git@github.com:git-town/git-town.git":                 {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("git")},
		"git@bitbucket.org/git-town/git-town.git":              {Host: "bitbucket.org", Org: "git-town", Repo: "git-town", User: Some("git")},
		"git@bitbucket.org/git-town/git-town.github.com":       {Host: "bitbucket.org", Org: "git-town", Repo: "git-town.github.com", User: Some("git")},
		"git@github.com:git-town/git-town":                     {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("git")},
		"git@gitlab.com:gitlab-com/www-gitlab-com.git":         {Host: "gitlab.com", Org: "gitlab-com", Repo: "www-gitlab-com", User: Some("git")},
		"git@gitlab.com:gitlab-com/www-gitlab-com":             {Host: "gitlab.com", Org: "gitlab-com", Repo: "www-gitlab-com", User: Some("git")},
		"git@gitlab.com:gitlab-org/quality/triage-ops.git":     {Host: "gitlab.com", Org: "gitlab-org/quality", Repo: "triage-ops", User: Some("git")},
		"git@gitlab.com:gitlab-org/quality/triage-ops":         {Host: "gitlab.com", Org: "gitlab-org/quality", Repo: "triage-ops", User: Some("git")},
		"git@git.example.com:4022/a/b.git":                     {Host: "git.example.com", Org: "a", Repo: "b", User: Some("git")},
		"git@git.example.com:4022/a/b":                         {Host: "git.example.com", Org: "a", Repo: "b", User: Some("git")},
		"username@bitbucket.org:git-town/git-town.git":         {Host: "bitbucket.org", Org: "git-town", Repo: "git-town", User: Some("username")},
		"username@bitbucket.org:git-town/git-town":             {Host: "bitbucket.org", Org: "git-town", Repo: "git-town", User: Some("username")},
		"https://github.com/git-town/git-town.git":             {Host: "github.com", Org: "git-town", Repo: "git-town", User: None[string]()},
		"https://github.com/git-town/git-town":                 {Host: "github.com", Org: "git-town", Repo: "git-town", User: None[string]()},
		"https://user@github.com/git-town/git-town.git":        {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("user")},
		"https://user@github.com/git-town/git-town":            {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("user")},
		"https://user:secret@github.com/git-town/git-town.git": {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("user:secret")},
		"https://user:secret@github.com/git-town/git-town":     {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("user:secret")},
		"http://github.com/git-town/git-town.git":              {Host: "github.com", Org: "git-town", Repo: "git-town", User: None[string]()},
		"http://github.com/git-town/git-town":                  {Host: "github.com", Org: "git-town", Repo: "git-town", User: None[string]()},
		"http://user@github.com/git-town/git-town.git":         {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("user")},
		"http://user@github.com/git-town/git-town":             {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("user")},
		"http://user:secret@github.com/git-town/git-town.git":  {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("user:secret")},
		"http://user:secret@github.com/git-town/git-town":      {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("user:secret")},
		"ssh://git@github.com/git-town/git-town.git":           {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("git")},
		"ssh://git@github.com/git-town/git-town":               {Host: "github.com", Org: "git-town", Repo: "git-town", User: Some("git")},
		"ssh://git@git.example.com:4022/a/b.git":               {Host: "git.example.com", Org: "a", Repo: "b", User: Some("git")},
		"ssh://git@git.example.com:4022/a/b":                   {Host: "git.example.com", Org: "a", Repo: "b", User: Some("git")},
		"git@ssh.dev.azure.com:v3/project/org/repo":            {Host: "ssh.dev.azure.com", Org: "project/org", Repo: "repo", User: Some("git")},
	}
	for give, want := range tests {
		have, has := giturl.Parse(give).Get()
		must.True(t, has)
		must.Eq(t, want, have)
	}
}
