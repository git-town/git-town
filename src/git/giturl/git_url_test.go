package giturl_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/giturl"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestParse(t *testing.T) {
	t.Parallel()
	tests := map[string]giturl.Parts{
		"git@github.com:git-town/git-town.git":                 {User: Some("git"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"git@bitbucket.org/git-town/git-town.git":              {User: Some("git"), Host: "bitbucket.org", Org: "git-town", Repo: "git-town"},
		"git@bitbucket.org/git-town/git-town.github.com":       {User: Some("git"), Host: "bitbucket.org", Org: "git-town", Repo: "git-town.github.com"},
		"git@github.com:git-town/git-town":                     {User: Some("git"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"git@gitlab.com:gitlab-com/www-gitlab-com.git":         {User: Some("git"), Host: "gitlab.com", Org: "gitlab-com", Repo: "www-gitlab-com"},
		"git@gitlab.com:gitlab-com/www-gitlab-com":             {User: Some("git"), Host: "gitlab.com", Org: "gitlab-com", Repo: "www-gitlab-com"},
		"git@gitlab.com:gitlab-org/quality/triage-ops.git":     {User: Some("git"), Host: "gitlab.com", Org: "gitlab-org/quality", Repo: "triage-ops"},
		"git@gitlab.com:gitlab-org/quality/triage-ops":         {User: Some("git"), Host: "gitlab.com", Org: "gitlab-org/quality", Repo: "triage-ops"},
		"git@git.example.com:4022/a/b.git":                     {User: Some("git"), Host: "git.example.com", Org: "a", Repo: "b"},
		"git@git.example.com:4022/a/b":                         {User: Some("git"), Host: "git.example.com", Org: "a", Repo: "b"},
		"username@bitbucket.org:git-town/git-town.git":         {User: Some("username"), Host: "bitbucket.org", Org: "git-town", Repo: "git-town"},
		"username@bitbucket.org:git-town/git-town":             {User: Some("username"), Host: "bitbucket.org", Org: "git-town", Repo: "git-town"},
		"https://github.com/git-town/git-town.git":             {User: None[string](), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"https://github.com/git-town/git-town":                 {User: None[string](), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"https://user@github.com/git-town/git-town.git":        {User: Some("user"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"https://user@github.com/git-town/git-town":            {User: Some("user"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"https://user:secret@github.com/git-town/git-town.git": {User: Some("user:secret"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"https://user:secret@github.com/git-town/git-town":     {User: Some("user:secret"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://github.com/git-town/git-town.git":              {User: None[string](), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://github.com/git-town/git-town":                  {User: None[string](), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://user@github.com/git-town/git-town.git":         {User: Some("user"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://user@github.com/git-town/git-town":             {User: Some("user"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://user:secret@github.com/git-town/git-town.git":  {User: Some("user:secret"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://user:secret@github.com/git-town/git-town":      {User: Some("user:secret"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"ssh://git@github.com/git-town/git-town.git":           {User: Some("git"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"ssh://git@github.com/git-town/git-town":               {User: Some("git"), Host: "github.com", Org: "git-town", Repo: "git-town"},
		"ssh://git@git.example.com:4022/a/b.git":               {User: Some("git"), Host: "git.example.com", Org: "a", Repo: "b"},
		"ssh://git@git.example.com:4022/a/b":                   {User: Some("git"), Host: "git.example.com", Org: "a", Repo: "b"},
	}
	for give, want := range tests {
		have, has := giturl.Parse(give).Get()
		must.True(t, has)
		must.Eq(t, want, have)
	}
}
